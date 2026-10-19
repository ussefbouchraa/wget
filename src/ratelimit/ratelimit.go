
package ratelimit

import (
	"io"
	"time"
)

const chunkSize = 32 * 1024 // 32 KiB

//-------------------reader
type Reader struct {
	src     io.Reader
	limiter *Limiter
}

func NewReader(src io.Reader, rate float64) *Reader {
	return &Reader{
		src:     src,
		limiter: newLimiter(rate),
	}
}

func (r *Reader) Read(p []byte) (int, error) {
	if len(p) > chunkSize {
		p = p[:chunkSize]
	}

	n, err := r.src.Read(p)
	if n > 0 {
		r.limiter.wait(n)
	}

	return n, err
}

//-----------------------Limiter

type Limiter struct {
	rate  float64
	read  int64
	start time.Time
}

func newLimiter(rate float64) *Limiter {
	return &Limiter{
		rate:  rate,
		start: time.Now(),
	}
}

func (l *Limiter) wait(n int) {
	l.read += int64(n)

	// Calculate when these bytes are allowed to be released.
	allowedAt := l.start.Add(
		time.Duration(float64(l.read) / l.rate * float64(time.Second)),
	)

	// Wait if the download is ahead of schedule.
	if delay := time.Until(allowedAt); delay > 0 {
		time.Sleep(delay)
	}
}
