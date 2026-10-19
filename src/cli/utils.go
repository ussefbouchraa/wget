
package cli

import (
	"fmt"
	"net/url"
	"strings"
)

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
