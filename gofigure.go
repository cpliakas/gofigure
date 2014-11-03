package gofigure

import (
	"flag"
	"os"
)

var flagSet *flag.FlagSet

type Config struct {
	EnvPrefix string
	options   map[string]*Option
	values    map[string]*string
}

func New() *Config {
	flagSet = flag.CommandLine
	return &Config{
		options: make(map[string]*Option),
		values:  make(map[string]*string),
	}
}

func (c *Config) Add(flag string) *Option {
	c.options[flag] = &Option{
		envVar: "",
		def:    "",
		desc:   "",
	}
	return c.options[flag]
}

func (c *Config) Get(name string) *string {
	return c.values[name]
}

// See https://github.com/rakyll/globalconf/blob/master/globalconf.go
func (c *Config) Parse() {

	for name, o := range c.options {
		c.values[name] = flag.String(name, "", o.desc)
	}

	// Map of whether command line args were passed
	passed := make(map[string]bool)
	flagSet.Visit(func(f *flag.Flag) {
		passed[f.Name] = true
	})

	flagSet.VisitAll(func(f *flag.Flag) {
		if passed[f.Name] {
			return
		}

		envVar := c.EnvPrefix + c.options[f.Name].envVar
		if val := os.Getenv(envVar); val != "" {
			flagSet.Set(f.Name, val)
		} else {
			flagSet.Set(f.Name, c.options[f.Name].def)
		}

	})

	flag.Parse()
}

type Option struct {
	envVar, def, desc string
}

func (o *Option) Default(def string) *Option {
	o.def = def
	return o
}

func (o *Option) EnvVar(envVar string) *Option {
	o.envVar = envVar
	return o
}

func (o *Option) Description(desc string) *Option {
	o.desc = desc
	return o
}
