package gofigure

import (
	"flag"
	"os"
)

// Pointer to the default set of command line flags.
var set *flag.FlagSet

// Config contains the configuration options that may be set by
// command line flags and environment variables.
type Config struct {
	EnvPrefix string
	options   map[string]*Option
	values    map[string]*string
}

// Returns a new Config instance.
func New() *Config {
	set = flag.CommandLine
	return &Config{
		options: make(map[string]*Option),
		values:  make(map[string]*string),
	}
}

// Adds a configuration option, returns an Option instance for
// easily setting the corresponding environment variable, default
// value, and description.
func (c *Config) Add(name string) *Option {
	c.options[name] = &Option{
		envVar: "",
		def:    "",
		desc:   "",
	}
	return c.options[name]
}

// Returns a configuration option by flag name.
func (c *Config) Get(name string) *string {
	return c.values[name]
}

// Parses the configuration options into defined flags, sets the value
// accordingly. Options are read first from command line flags, then
// from environment variables, and falls back to the default value if
// neither are set.
//
// See https://github.com/rakyll/globalconf/blob/master/globalconf.go
func (c *Config) Parse() {

	// Sets the flags from the configuration options.
	for name, o := range c.options {
		c.values[name] = flag.String(name, o.def, o.desc)
	}

	flag.Parse()

	// Gather the flags passed through command line.
	passed := make(map[string]bool)
	set.Visit(func(f *flag.Flag) {
		passed[f.Name] = true
	})

	set.VisitAll(func(f *flag.Flag) {

		// Skip flags passed through the command line as the option is
		// already set and takes precedence over environment variables.
		if passed[f.Name] {
			return
		}

		// Skip flags that aren't added to Config.
		if _, isset := c.options[f.Name]; !isset {
			return
		}

		// Some options shouldn't be set via environment variables.
		if c.options[f.Name].envVar == "" {
			return
		}

		// If the configuration option was not passed via the command line,
		// check the corresponding environment variable.
		envVar := c.EnvPrefix + c.options[f.Name].envVar
		if val := os.Getenv(envVar); val != "" {
			set.Set(f.Name, val)
		}
	})
}

// Option contains the details of a configuration options,
// e.g. corresponding environment variable, default value,
// description.
type Option struct {
	envVar, def, desc string
}

// Sets the configuration option's default value.
func (o *Option) Default(def string) *Option {
	o.def = def
	return o
}

// Sets the configuration option's corresponding environment variable.
func (o *Option) EnvVar(envVar string) *Option {
	o.envVar = envVar
	return o
}

// Sets the configuration options long description.
func (o *Option) Description(desc string) *Option {
	o.desc = desc
	return o
}
