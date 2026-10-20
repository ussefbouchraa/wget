// It parses CLI flags, validates arguments and dispatches to the right
// download mode (single file, batch, background, mirror).

package main

import (
	"fmt"
	"os"
	"time"

	"wget/src/batch"
	"wget/src/cli"
	"wget/src/downloader"
	"wget/src/logging"
)

func main() {
	conf, err := cli.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "wget:", err)
		os.Exit(1)
	}

	// Background mode: print the notice, then send everything to wget-log.
	var logFile *os.File
	if conf.Background {
		logFile, err = logging.Start()
		if err != nil {
			fmt.Fprintln(os.Stderr, "wget:", err)
			os.Exit(1)
		}
		defer logging.Stop(logFile)
	}

	start := time.Now()
	fmt.Printf("start at %s\n", start.Format("2006-01-02 15:04:05"))

	if conf.InputFile != "" {
		err = batch.Run(conf)
	} else {
		err = downloader.Download(conf, os.Stdout)
	}

	fmt.Printf("finished at %s\n", time.Now().Format("2006-01-02 15:04:05"))

	if err != nil {
		fmt.Fprintln(os.Stderr, "wget:", err)
		os.Exit(1)
	}
}
