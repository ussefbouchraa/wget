package progress

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

const barWidth = 112 // progress width

// Bar shows how much of a download is complete.
type Bar struct {
	start    time.Time
	fileSize int64 //total bytes
	received int64 // bytes downloaded

	out    io.Writer
	mu     sync.Mutex
	stopCh chan struct{}
	doneCh chan struct{}
}

// NewBar starts a bar that redraws itself every 100 ms on out.
// A nil out means os.Stdout.
func NewBar(fileSize int64, out io.Writer) *Bar {
	if out == nil {
		out = os.Stdout
	}
	b := &Bar{fileSize: fileSize, out: out, start: time.Now(), stopCh: make(chan struct{}), doneCh: make(chan struct{})}
	go b.refresh()
	return b
}

// Write counts the bytes flowing through the bar, so a Bar can be
// wired into io.TeeReader.
func (b *Bar) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.received += int64(len(p))
	return len(p), nil
}

// Finish draws the last frame and moves to the next line.
func (b *Bar) Finish() {
	close(b.stopCh)
	<-b.doneCh
	b.render()
	fmt.Fprintln(b.out)
}

func (b *Bar) refresh() {
	defer close(b.doneCh)
	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-b.stopCh:
			return
		case <-t.C:
			b.render()
		}
	}
}

// draw prints the current state on top of the same line.
func (b *Bar) render() {
	b.mu.Lock()
	defer b.mu.Unlock()

	received, elapsed := b.received, time.Since(b.start)
	if elapsed < time.Millisecond {
		elapsed = time.Millisecond
	}

	pct, filled, fileSize := 100.0, barWidth, "?"
	if b.fileSize > 0 {
		if received > b.fileSize {
			received = b.fileSize
		}
		pct = 100 * float64(received) / float64(b.fileSize)
		filled = int(barWidth * float64(received) / float64(b.fileSize))
		fileSize = formatSize(b.fileSize)
	} else if received == 0 {
		pct, filled = 0, 0
	}

	bar := strings.Repeat("=", filled) + strings.Repeat(" ", barWidth-filled)
	speed := formatSize(int64(float64(received) / elapsed.Seconds()))
	fmt.Fprintf(b.out, "\r %s / %s [%s] %.2f%% %s/s %s\x1b[K",
		formatSize(received), fileSize, bar, pct, speed, timeRemaining(pct, elapsed))
}

// size formats bytes with binary units, e.g. "55.05 KiB".
func formatSize(n int64) string {
	if n < 0 {
		return "?"
	}
	units := []string{"B", "KiB", "MiB", "GiB", "TiB"}
	i, f := 0, float64(n)
	for f >= 1024 && i < len(units)-1 {
		f /= 1024
		i++
	}
	if i == 0 {
		return fmt.Sprintf("%d B", n)
	}
	return fmt.Sprintf("%.2f %s", f, units[i])
}

// timeRemaining estimates the time left: "--" at 0%, "0s" when complete.
func timeRemaining(pct float64, elapsed time.Duration) string {
	if pct <= 0 {
		return "--"
	}
	if pct >= 100 {
		return "0s"
	}
	secs := int(elapsed.Seconds() * (100 - pct) / pct)
	return fmt.Sprintf("%dm%ds", secs/60, secs%60)
}
