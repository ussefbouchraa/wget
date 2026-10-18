// It parses CLI flags, validates arguments and dispatches to the right
// download mode (single file, batch, background, mirror).

package main

import (
	"fmt"
	"os"
	"time"

	"wget/src/cli"
	"wget/src/downloader"
)

func main() {
	opts, err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "wget:", err)
		os.Exit(1)
	}

	start := time.Now()
	fmt.Printf("start at %s\n", start.Format("2006-01-02 15:04:05"))

	if err := downloader.Download(opts); err != nil {
		fmt.Fprintln(os.Stderr, "wget:", err)
		os.Exit(1)
	}

	fmt.Printf("finished at %s\n", time.Now().Format("2006-01-02 15:04:05"))
}
