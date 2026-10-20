// Package batch implements the -i flag: it reads a file of URLs and
// downloads them all asynchronously, one goroutine per URL.
package batch

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"

	"wget/src/cli"
	"wget/src/downloader"
)

// printMu keeps each download's output block atomic on the shared stream.
var printMu sync.Mutex

// Run downloads every URL listed in conf.InputFile, one goroutine each.
// Each download buffers its own output and prints it in one shot when it finish
func Run(conf *cli.Wconfig) error {
	urls, err := readURLs(conf.InputFile)
	if err != nil {
		return err
	}
	if len(urls) == 0 {
		return fmt.Errorf("no urls found in %s", conf.InputFile)
	}

	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()

			c := *conf // copy the shared settings
			c.URL = u
			c.OutputName = ""
			c.InputFile = ""
			c.Background = false

			var buf bytes.Buffer
			if err := downloader.Download(&c, &buf); err != nil {
				fmt.Fprintf(&buf, "wget: %v\n", err)
			}

			printMu.Lock()
			os.Stdout.Write(buf.Bytes())
			printMu.Unlock()
		}(u)
	}
	wg.Wait()

	fmt.Println("All downloads completed.")
	return nil
}

// readURLs returns the non-empty trimmed lines of the input file.
func readURLs(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open input file %s: %w", path, err)
	}
	defer f.Close()

	var urls []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		urls = append(urls, line)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("cannot read input file %s: %w", path, err)
	}
	return urls, nil
}
