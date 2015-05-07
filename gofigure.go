package gofigure

import (
	"fmt"
	"github.com/droundy/goopt"
	"log"
	"os"
)

// This function is a convenience function. It displays the default values for us but handles them
// in the way we prefer it.
func GooptFigureString(names []string, def string, help string) *string {
	s := new(string)
	*s = ""
	f := func (ss string) error {
		*s = ss
		return nil
	}
	goopt.ReqArg(names, def, help, f)
	return s
}

// Config contains the configuration options that may be set by
// command line flags and environment variables.
type Config struct {
	Description			string
	DescribeEnvironment bool	// if true, the environment variable is automatically added to the flag description
	DisableCommandLine	bool
	EnvPrefix			string
	EnvOverridesFile	bool
	FileParser			File
	RequireFile			bool    // application will panic if RequireFile == true, FileParser != nil and file doesn't exist
	Version				string
	options				map[string]*Option
	flags				map[string]*string
	values				map[string]*string
}

// Returns a new Config instance.
func New() *Config {
	return &Config{
		DescribeEnvironment: false,
		DisableCommandLine:	 false,
		EnvOverridesFile:    false,
		RequireFile:         true,
		options:			 make(map[string]*Option),
		flags:				 make(map[string]*string),
		values:				 make(map[string]*string),
	}
}

// Adds a configuration option, returns an Option instance for
// easily setting the corresponding environment variable, default
// value, and description.
func (c *Config) Add(name string) *Option {
	c.options[name] = &Option{
		name:    name,
		envVar:  "",
		def:     "",
		desc:    "",
		longOpt: name,
	}
	return c.options[name]
}

// Returns a configuration option by flag name.
func (c Config) Get(name string) *string {
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
		cmdline = append(cmdline, "--"+o.longOpt)
		desc := o.desc
		if c.DescribeEnvironment && o.envVar != "" {
			desc += fmt.Sprintf(" Environment variable: %s%s.", c.EnvPrefix, o.envVar)
		}
		c.flags[name] = GooptFigureString(cmdline, o.def, desc)
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
		err := c.ParseFile(passed)
		if err != nil && c.RequireFile {
			log.Panicf("File defined but could not be parsed: %s", err.Error())
		} else if err != nil {
			log.Printf("Could not parse file: %s", err.Error())
		}
	} else {
		err := c.ParseFile(passed)
		if err != nil && c.RequireFile {
			log.Panicf("File defined but could not be parsed: %s", err.Error())
		} else if err != nil {
			log.Printf("Could not parse file: %s", err.Error())
		}
		c.ParseEnv(passed)
	}
}

func (c *Config) ParseEnv(passed map[string]bool) {
	for name, f := range c.options {

		// Skip flags passed through the command line as the option is
		// already set and takes precedence over environment variables.
		if passed[name] {
			continue
		}

		// Some options shouldn't be set via environment variables.
		if f.envVar == "" {
			continue
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

func (c Config) parseFileToMap(handler File) (ValueMap, error) {
	root, err := handler.Parse()
	if err != nil {
		return nil, err
	}
	ret := ValueMap{}
	for _, opt := range c.options {
		if val, ok := root.FindOption(opt); ok {
			ret.Set(opt.GetName(), val)
		}
	}
	return ret, nil
}

func (c *Config) ParseFile(passed map[string]bool) error {
	if c.FileParser == nil {
		return nil
	}

	values, err := c.parseFileToMap(c.FileParser)
	if err != nil {
		return err
	}

	for k, v := range values {
		if passed[k] {
			continue
		}

		if p, ok := c.values[k]; ok {
			*p = v
			passed[k] = true
		}
	}
	return nil
}

// Returns all options.
func (c *Config) GetOptions() map[string]*Option {
	ret := make(map[string]*Option)
	for name, val := range c.options {
		ret[name] = val
	}
	return ret
}

// Option contains the details of a configuration options,
// e.g. corresponding environment variable, default value,
// description.
type Option struct {
	name, envVar, shortOpt, def, desc, longOpt string
	fileSpec                                   string // The file spec is of the form "(CATEGORY.)*NAME", eg. for 'foo' under the category 'bar', it would be foo.bar
}

func (o Option) GetName() string {
	return o.name
}

func (o Option) GetFileSpec() string {
	return o.fileSpec
}

func (o Option) GetShortOpt() string {
	return o.shortOpt
}

func (o Option) GetLongOpt() string {
	return o.longOpt
}

func (o Option) GetDefault() string {
	return o.def
}

func (o Option) GetEnvVar() string {
	return o.envVar
}

func (o Option) GetDescription() string {
	return o.desc
}

func (o *Option) FileSpec(spec string) *Option {
	o.fileSpec = spec
	return o
}

func (o *Option) ShortOpt(opt string) *Option {
	o.shortOpt = opt
	return o
}

func (o *Option) LongOpt(opt string) *Option {
	o.longOpt = opt
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
