
package logging

import (
	"fmt"
	"os"
)

// LogFile is the name of the background log, like GNU wget.
const LogFile = "wget-log"

// savedOut and savedErr hold the real terminal streams while redirected.
var (
	savedOut *os.File
	savedErr *os.File
)

// Start prints the notice, then redirects stdout and stderr to the log file.
func Start() (*os.File, error) {
	fmt.Printf("Output will be written to '%s'.\n", LogFile)

	f, err := os.Create(LogFile)
	if err != nil {
		return nil, fmt.Errorf("cannot create %s: %w", LogFile, err)
	}

	savedOut, savedErr = os.Stdout, os.Stderr
	os.Stdout, os.Stderr = f, f

	return f, nil
}

// Stop restores the terminal streams and closes the log file.
func Stop(f *os.File) {
	os.Stdout, os.Stderr = savedOut, savedErr
	f.Close()
}
