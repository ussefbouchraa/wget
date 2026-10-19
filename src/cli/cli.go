// Package cli handles command-line arguments and flags.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Wconfig stores the download settings.
type Wconfig struct {
	URL        string
	OutputName string
	Path       string
	RateLimit  float64 // bytes per second, 0 means "no limit"
}

// Parse reads the arguments and checks if they are valid.
func Parse(args []string) (*Wconfig, error) {
	fs := flag.NewFlagSet("wget", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var conf Wconfig
	var rateLimit string
	fs.StringVar(&conf.OutputName, "O", "", "output file name")
	fs.StringVar(&conf.Path, "P", "", "output directory")
	fs.StringVar(&rateLimit, "rate-limit", "", "max download speed, e.g. 300k or 2M")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if rateLimit != "" {
		bps, err := ParseRateLimit(rateLimit)
		if err != nil {
			return nil, err
		}
		conf.RateLimit = bps //speed limit
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

// ParseRateLimit turns a value like "300k", "2M" or "512" into anumber of bytes per second	
func ParseRateLimit(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, errors.New("rate limit value cannot be empty")
	}

	multiplier := 1.0
	switch s[len(s)-1] {
	case 'k', 'K':
		multiplier = 1024
		s = s[:len(s)-1]
	case 'm', 'M':
		multiplier = 1024 * 1024
		s = s[:len(s)-1]
	case 'g', 'G':
		multiplier = 1024 * 1024 * 1024
		s = s[:len(s)-1]
	}

	value, err := strconv.ParseFloat(s, 64)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("invalid rate limit %q (use e.g. 300k or 2M)", s)
	}
	return value * multiplier, nil
}
