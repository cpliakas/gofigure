package gofigure

import (
	"github.com/droundy/goopt"
	"os"
)

// Config contains the configuration options that may be set by
// command line flags and environment variables.
type Config struct {
	Description			string
	DisableCommandLine	bool
	EnvPrefix			string
	EnvOverridesFile	bool
	Filename			string
	FileParser			File
	Version				string
	options				map[string]*Option
	flags				map[string]*string
	values				map[string]*string
}

// Returns a new Config instance.
func New() *Config {
	return &Config{
		DisableCommandLine:	false,
		EnvOverridesFile:   false,
		options:			make(map[string]*Option),
		flags:				make(map[string]*string),
		values:				make(map[string]*string),
	}
}

// Adds a configuration option, returns an Option instance for
// easily setting the corresponding environment variable, default
// value, and description.
func (c *Config) Add(name string) *Option {
	c.options[name] = &Option{
		name:   name,
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
	goopt.Description = func() string {
		return c.Description 
	}

	goopt.Version = c.Version

	// Sets the flags from the configuration options.
	for name, o := range c.options {
		cmdline := []string{}
		if o.shortOpt != "" {
			cmdline = append(cmdline, "-" + o.shortOpt)
		}
		cmdline = append(cmdline, "--"+name)
		c.flags[name] = goopt.String(cmdline, "", o.desc)
		defcopy := o.def
		c.values[name] = &defcopy
	}

	passed := make(map[string]bool)

	if !c.DisableCommandLine {
		goopt.Parse(nil)

		// Gather the options passed through command line.
		for name, f := range c.flags {
			if *f != "" {
				passed[name] = true
				*c.values[name] = *f
			}
		}
	}

	if (c.EnvOverridesFile) {
		c.ParseEnv(passed)
		c.ParseFile(passed)
	} else {
		c.ParseFile(passed)
		c.ParseEnv(passed)
	}
}

func (c *Config) ParseEnv(passed map[string]bool) {
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
			*c.values[name] = val
			passed[name] = true
		}
	}
}

func (c Config) parseFileToMap(name string, handler File) (ValueMap, error) {
	root, err := handler.Parse(name)
	if err != nil {
		return nil, err
	}
	ret := ValueMap{}
	for _, opt := range c.options {
		if val, ok := root.FindOption(opt); ok {
			ret.Set(opt.Name(), val)
		}
	}
	return ret, nil
}

func (c *Config) ParseFile(passed map[string]bool) error {
	if (c.FileParser == nil || c.Filename == "") {
		return nil
	}

	values, err := c.parseFileToMap(c.Filename, c.FileParser)
	if err != nil {
		return err
	}

	for k, v := range values {
		if p, ok := c.values[k]; ok {
			*p = v
			passed[k] = true
		}
	}
	return nil
}

// Option contains the details of a configuration options,
// e.g. corresponding environment variable, default value,
// description.
type Option struct {
	name, envVar, shortOpt, def, desc string
	fileSpec				string              // The file spec is of the form "(CATEGORY.)*NAME", eg. for 'foo' under the category 'bar', it would be foo.bar
}

func (o Option) Name() string {
	return o.name
}

func (o *Option) ShortOpt(opt string) *Option {
	o.shortOpt = opt
	return o
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
