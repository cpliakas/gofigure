package gofigure

import (
	"github.com/droundy/goopt"
	"os"
)

// Config contains the configuration options that may be set by
// command line flags and environment variables.
type Config struct {
	Description	string
	Version		string
	EnvPrefix	string
	options		map[string]*Option
	flags		map[string]*string
	values		map[string]string
}

// Returns a new Config instance.
func New() *Config {
	return &Config{
		options: make(map[string]*Option),
		flags:  make(map[string]*string),
		values: make(map[string]string),
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
func (c *Config) Get(name string) string {
	return c.values[name]
}

// Parses the configuration options into defined flags, sets the value
// accordingly. Options are read first from command line flags, then
// from environment variables, and falls back to the default value if
// neither are set.
//
// See https://github.com/rakyll/globalconf/blob/master/globalconf.go
func (c *Config) Parse() {
	goopt.Description = func() string {
		return c.Description 
	}

	goopt.Version = c.Version

	// Sets the flags from the configuration options.
	for name, o := range c.options {
		c.flags[name] = goopt.String([]string{"--"+name}, "", o.desc)
		c.values[name] = o.def
	}

	goopt.Parse(nil)

	// Gather the flags passed through command line.
	passed := make(map[string]bool)
	for name, f := range c.flags {
		if *f != "" {
			passed[name] = true
			c.values[name] = *f
		}
	}

	for name, f := range c.options {

		// Skip flags passed through the command line as the option is
		// already set and takes precedence over environment variables.
		if passed[name] {
			return
		}

		// Some options shouldn't be set via environment variables.
		if f.envVar == "" {
			return
		}

		// If the configuration option was not passed via the command line,
		// check the corresponding environment variable.
		envVar := c.EnvPrefix + f.envVar
		if val := os.Getenv(envVar); val != "" {
			c.values[name] = val
		}
	}
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
