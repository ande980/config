# Config

I needed a lite configuration library for a personal Go project, and having read Peter Bourgon's
[treatise](https://peter.bourgon.org/go-for-industrial-programming/) on industrial programming I 
thought I would have an opinionated attempt at it.

The library allows the use of command line flags, environmental variables, and / or JSON encoded
text files. Flags and env vars are generated automatically from struct property names and tags. The
standard library `encoding/json` encoder is used for JSON files. JSON sucks for configuration 
(comments) but I wanted to keep the import graph as small a possible initially. It would be trivial 
to allow toml (preferred) or yaml (please don't) encoded files.

At a very basic level the rationale was to create a library that would allow an instance of a struct
to be passed to it at runtime and layered providers would supply configuration values from the 
environment.

## Assumptions
* The library trades startup time for ease of use which in a long running API shouldn't cause too 
many problems for most people, unless you're running in the cloud.
* Flags are created in their own flagset because `flag.ExitOnError` is the default and I don't want
that. This also means that if you have other flags specified on the command line (perhaps you want
to interpret them with another library) then the flag library will stop processing and return an 
error. I'm assuming that if you have a configuration struct then all of the configuration is in it 
and you want this library to handle it.
* This library does not panic in an of itself. However, although care has been taken not to cause a
panic from the reflect library I can't guarantee I found them all.

## Usage
Most common usage would simply be this:
```go
package main

import (
    "fmt"
    "os"
    
    "github.com/ande980/config"
)

func main() {
    if err := run(); err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        os.Exit(2)
    }
}

func run() error {
    cfg := &Cfg{}
    if err := config.Parse(cfg); err != nil {
        if err == flag.ErrHelp {
            return nil
        }
        return err    
    }
    // use cfg
}
```

The library includes two optional interfaces that can be implemented:
1) `Initer` just in case you can't, for some reason, create a constructor for your type or have some 
theoretical requirement to set default values outside of a constructor ... or something.
2) `Validator` if you want to use JSON schema or some custom business logic to ensure your configuration 
is valid.

Using the `Validator` interface would look something like this:
```go
package main

import (
    "fmt"
    "os"
    
    "github.com/ande980/config"
)

type Cfg struct {
    Addr      string `usage:"Local address to bind to"`
    ProxyAddr string `usage:"Remote proxy address"`
}

func(c *Cfg) Validate() error {
    if c.Addr == "" {
        c.Addr = ":0"
    }

    if c.ProxyAddr == "" {
        return fmt.Errorf("validate: proxy address must be provided")
    }
    return nil
}

func main() {
    if err := run(); err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        os.Exit(2)
    }
}

func run() error {
    cfg := &Cfg{}
    if err := config.Parse(cfg); err != nil {
        if err == flag.ErrHelp {
            return nil
        }
        return err    
    }
    // use cfg
}
```

## Additional Providers
Both a toml, and a yaml, provider have been created but neither are registered by default.
This is to keep the import graph small and restricted to (almost) only standard libraries.
Only one file provider can be registered at once, so registering the toml provider for example
will unregister the json or yaml provider if either are registered.

## TODO
- [x] Either add in panic recovery or change reflection panics to errors  
- [x] Add toml support
- [x] Add yaml support
- [x] More canonical flag names
- [x] Tests for toml and yaml providers  
- [ ] Generate default usage verbiage if none is provided  