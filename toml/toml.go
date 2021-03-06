package toml

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Provider is a config provider that reads from a toml
// file or io.Reader and scans into the specified struct.
type Provider struct {
	r   io.Reader
	err error
}

// New is the default way to create a toml Provider. The entire
// file is read when New is called and a reader created from
// the buffer. If there is an error reading the file it
// is stored in Provider and returned during Parse.
func New() *Provider {
	filepathNoExt := strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))
	return WithPath(filepathNoExt + ".toml")
}

// WithPath allows for a non-standard configuration file to be
// specified at runtime. If there is an error reading the file it
// is stored in Provider and returned during Parse.
func WithPath(path string) *Provider {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Provider{strings.NewReader(""), nil} // No-op reader, but one that doesn't generate io.EOF
	}

	buf, err := ioutil.ReadFile(path)
	p := &Provider{r: bytes.NewReader(buf)}
	if err != nil {
		p.err = fmt.Errorf("reading configuration file: %v", err)
	}
	return p
}

// WithReader accepts a reader and returns a toml Provider.
func WithReader(r io.Reader) *Provider {
	return &Provider{r, nil}
}

// Parse implements the config.Provider interface.
func (p *Provider) Parse(i interface{}) error {
	if p.err != nil {
		return p.err
	}

	if _, err := toml.DecodeReader(p.r, i); err != nil && err != io.EOF {
		return fmt.Errorf("decoding toml file: %v", err)
	}

	return nil
}
