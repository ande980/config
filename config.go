package config

import (
	"config/env"
	"config/flags"
	"config/json"
	"reflect"
	"sync"

	multierror "github.com/hashicorp/go-multierror"
)

func init() {
	Register(env.New())
	Register(json.New())
	Register(flags.New())
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

var (
	mu        sync.Mutex
	once      sync.Once
	providers []Provider
)

// Register accepts a provider and registeres it for runtime.
// Unlike sql driver registrations, multiple registrations per
// type are allowed but each successive call will replace the last.
// This is to allow a default to be registered, but a more specific
// variant to be registered later in place of it.
func Register(p Provider) {
	once.Do(func() {
		providers = []Provider{}
	})

	mu.Lock()
	providers = append(providers, p)
	mu.Unlock()
}

// Parse acccepts a variadic number of config providers and returns an error.
// If a single provider returns an error then it will be return even if
// all other providers functioned correctly.
func Parse(i interface{}) error {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Ptr {
		panic(&reflect.ValueError{Method: "config.Register", Kind: reflect.Ptr})
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		panic(&reflect.ValueError{Method: "config.Register", Kind: reflect.Struct})
	}

	if len(providers) == 0 {
		panic("no registered providers")
	}

	if initer, ok := i.(Initer); ok {
		if err := initer.Init(); err != nil {
			panic(err)
		}
	}

	var result *multierror.Error
	mu.Lock()
	for _, provider := range providers {
		if err := provider.Parse(i); err != nil {
			result = multierror.Append(result, err)
		}
	}
	mu.Unlock()

	if validator, ok := i.(Validator); ok {
		if err := validator.Validate(); err != nil {
			panic(err)
		}
	}

	return result.ErrorOrNil()
}
