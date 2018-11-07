package config

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/ande980/config/env"
	"github.com/ande980/config/flags"
	"github.com/ande980/config/json"
	"github.com/ande980/config/toml"
	"github.com/ande980/config/yaml"
	multierror "github.com/hashicorp/go-multierror"
)

var (
	// ErrHelp is returned when the -h or --help flags are used.
	ErrHelp = errors.New("help requested")
	// ErrVersion is returned when the -v or --version flags are used.
	ErrVersion = errors.New("version requested")
)

// Yay for global state. Why are you parsing more than one configuration file?
var providers = []Provider{
	env.New(),
	flags.New(),
}

// Initer is an optional interface that configuration structs
// can implement with a single method that will be called
// before any values are scanned into the struct. This is
// useful as a constructor without having a lot of
// initialization code at the call site.
type Initer interface {
	Init() error
}

// Validator is an optional interface that configuration structs
// can implement with a single method that will be called after
// scanning has completed. This is useful as a validation phase.
type Validator interface {
	Validate() error
}

// ProviderFunc is a convenience type a la http.HandlerFunc.
type ProviderFunc func(interface{}) error

// Parse allows ProviderFunc to satisfy the Provider interface
func (p ProviderFunc) Parse(i interface{}) error {
	return p(i)
}

// Provider describes a single function to take an arbitrary struct and infill.
type Provider interface {
	Parse(interface{}) error
}

// Parse acccepts a variadic number of config providers and returns an error.
// If a single provider returns an error then it will be return even if
// all other providers functioned correctly.
func Parse(i interface{}) (err error) {
	defer func() {
		if p := recover(); p != nil {
			switch t := p.(type) {
			case error:
				err = t
			default:
				err = fmt.Errorf("%v", t)
			}
		}
	}()

	// This is highly opinionated but it does what I need it to.
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Usage = func() {}
	fs.SetOutput(ioutil.Discard)
	fs.Parse(os.Args[1:])
	if len(fs.Args()) > 0 {
		configPath := fs.Args()[0]
		if _, err := os.Stat(configPath); err == nil {
			switch filepath.Ext(configPath) {
			case ".json":
				providers = append(providers, json.WithPath(configPath))
			case ".toml":
				providers = append(providers, toml.WithPath(configPath))
			case ".yaml", ".yml":
				providers = append(providers, yaml.WithPath(configPath))
			}
		} else {
			providers = append(providers, json.New(), toml.New(), yaml.New())
		}
	} else {
		providers = append(providers, json.New(), toml.New(), yaml.New())
	}

	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		err = &reflect.ValueError{Method: "parser.Parse", Kind: reflect.Ptr}
		return err
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		err = &reflect.ValueError{Method: "parser.Parse", Kind: reflect.Struct}
		return err
	}

	if len(providers) == 0 {
		err = fmt.Errorf("no providers specified")
		return err
	}

	if initer, ok := i.(Initer); ok {
		if err = initer.Init(); err != nil {
			return err
		}
	}

	var result *multierror.Error
	for _, provider := range providers {
		if err = provider.Parse(i); err != nil {
			switch err {
			case flags.ErrHelp:
				return ErrHelp
			case flags.ErrVersion:
				return ErrVersion
			default:
				result = multierror.Append(result, err)
			}
		}
	}

	if validator, ok := i.(Validator); ok {
		if err = validator.Validate(); err != nil {
			return err
		}
	}

	err = result.ErrorOrNil()
	return err
}
