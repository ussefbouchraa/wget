// It sends the GET request, checks the response status, reports the
// content length and streams the body to the destination file.

package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"wget/src/cli"
	"wget/src/progress"
	"wget/src/ratelimit"
)

// Download gets a file from the URL and saves it.
func Download(conf *cli.Wconfig) error {
	resp, err := http.Get(conf.URL)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	fmt.Printf("sending request, awaiting response... status %s\n", resp.Status)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server responded with %s", resp.Status)
	}

	size := resp.ContentLength
	fmt.Printf("content size: %d [%s]\n", size, humanSize(size))

	dest := conf.FullPath()
	fmt.Printf("saving file to: %s\n", dest)

	// Create the folder if it doesn't exist.
	if dir := filepath.Dir(dest); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("cannot create directory %s: %w", dir, err)
		}
	}

	// Create the file.
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("cannot create file %s: %w", dest, err)
	}
	defer f.Close()

	// When a rate limit is set, data flows through the limiter
	// first so it never arrives faster than allowed.
	var body io.Reader = resp.Body
	if conf.RateLimit > 0 {
		body = ratelimit.NewReader(resp.Body, conf.RateLimit)
	}

	// Stream the data into the file while updating the progress bar.
	bar := progress.NewBar(size)
	_, err = io.Copy(f, io.TeeReader(body, bar))
	bar.Finish()
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	fmt.Println()
	fmt.Printf("Downloaded [%s]\n", conf.URL)
	return nil
}

// humanSize converts bytes into a readable size like KB or MB.
func humanSize(n int64) string {
	if n < 0 {
		return "unknown"
	}
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%dB", n)
	}
	div, exp := int64(unit), 0
	for m := n / unit; m >= unit; m /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("~%.2f%cB", float64(n)/float64(div), "KMGTPE"[exp])
}
