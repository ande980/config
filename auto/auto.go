package auto

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/ande980/config"
	"github.com/ande980/config/json"
	"github.com/ande980/config/toml"
	"github.com/ande980/config/yaml"
)

// This is very icky but its one of the few ways I can think of to:
// 1: allow for picking up of an unnamed parameter as a config file
// 2: keep the default import graph limited to std libs + multierror
func init() {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.Parse(os.Args[1:])
	if len(fs.Args()) > 0 {
		configPath := fs.Args()[0]
		if _, err := os.Stat(configPath); err == nil {
			switch filepath.Ext(configPath) {
			case ".json":
				config.Register(json.WithPath(configPath))
			case ".toml":
				config.Register(toml.WithPath(configPath))
			case ".yaml", ".yml":
				config.Register(yaml.WithPath(configPath))
			}
		}
	}
}
