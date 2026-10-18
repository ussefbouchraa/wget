// Package cli handles command-line arguments and flags.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

// Wconfig stores the download settings.
type Wconfig struct {
	URL        string
	OutputName string
	Path       string
}

// Parse reads the arguments and checks if they are valid.
func Parse(args []string) (*Wconfig, error) {
	fs := flag.NewFlagSet("wget", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var conf Wconfig
	fs.StringVar(&conf.OutputName, "O", "", "output file name")
	fs.StringVar(&conf.Path, "P", "", "output directory")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if fs.NArg() != 1 {
		return nil, errors.New("usage: wget [flags] <url>")
	}

	conf.URL = fs.Arg(0)
	if err := validateURL(conf.URL); err != nil {
		return nil, err
	}

	return &conf, nil
}

// validateURL checks that the URL is valid and uses HTTP or HTTPS.
func validateURL(raw string) error {
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid url %q: %w", raw, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported scheme %q in url %q (only http/https)", u.Scheme, raw)
	}
	if u.Host == "" {
		return fmt.Errorf("missing host in url %q", raw)
	}
	return nil
}

// FileName returns the name used to save the file.
func (conf *Wconfig) FileName() string {
	if conf.OutputName != "" {
		return conf.OutputName
	}
	u, err := url.Parse(conf.URL)
	if err != nil {
		return "index.html"
	}
	name := u.Path
	if i := strings.LastIndex(name, "/"); i >= 0 {
		name = name[i+1:]
	}
	if name == "" {
		return "index.html"
	}
	return name
}

// FullPath returns the directory and file name together.
func (conf *Wconfig) FullPath() string {
	name := conf.FileName()
	if conf.Path == "" {
		return "./" + name
	}
	p := strings.TrimRight(conf.Path, "/")
	return p + "/" + name
}
