package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Provider is a config provider that reads from a JSON
// file or io.Reader and scans into the specified struct.
type Provider struct {
	r   io.Reader
	err error
}

// New is the default way to create a json Provider. The entire
// file is read when New is called and a reader created from
// the buffer. If there is an error reading the file it
// is stored in Provider and returned during Parse.
func New() *Provider {
	filepathNoExt := strings.TrimSuffix(filepath.Base(os.Args[0]), filepath.Ext(os.Args[0]))
	return WithPath(filepathNoExt + ".json")
}

// WithPath allows for a non-standard configuration file to be
// specified at runtime. If there is an error reading the file it
// is stored in Provider and returned during Parse.
func WithPath(path string) *Provider {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &Provider{strings.NewReader("{}"), nil} // No-op reader, but one that doesn't generate io.EOF
	}

	buf, err := ioutil.ReadFile(path)
	p := &Provider{r: bytes.NewReader(buf)}
	if err != nil {
		p.err = fmt.Errorf("reading configuration file: %v", err)
	}
	return p
}

// WithReader accepts a reader and returns a json Provider.
func WithReader(r io.Reader) *Provider {
	return &Provider{r, nil}
}

// Parse implements the config.Provider interface.
func (p *Provider) Parse(i interface{}) error {
	if p.err != nil {
		return p.err
	}

	if err := json.NewDecoder(p.r).Decode(i); err != nil && err != io.EOF {
		return fmt.Errorf("decoding json file: %v", err)
	}
	return nil
}
