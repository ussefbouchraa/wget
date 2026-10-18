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

// Options stores the download settings.
type Options struct {
	URL        string
	OutputName string
	Path       string
}

// Parse reads the arguments and checks if they are valid.
func Parse(args []string) (*Options, error) {
	fs := flag.NewFlagSet("wget", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var opts Options
	fs.StringVar(&opts.OutputName, "O", "", "output file name")
	fs.StringVar(&opts.Path, "P", "", "output directory")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if fs.NArg() != 1 {
		return nil, errors.New("usage: wget [flags] <url>")
	}

	opts.URL = fs.Arg(0)
	if err := validateURL(opts.URL); err != nil {
		return nil, err
	}

	return &opts, nil
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
func (o *Options) FileName() string {
	if o.OutputName != "" {
		return o.OutputName
	}
	u, err := url.Parse(o.URL)
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
func (o *Options) FullPath() string {
	name := o.FileName()
	if o.Path == "" {
		return "./" + name
	}
	p := strings.TrimRight(o.Path, "/")
	return p + "/" + name
}
