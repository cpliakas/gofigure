# gofigure

A configuration utility for Go insipred by [The Twelve-Factor App](http://12factor.net/config).

Configuration options are defined in command line arguments and
environment variables, evaluated in that order.

## Usage

```
import "github.com/cpliakas/gofigure"
```

```go
config := gofigure.New()
config.envPrefix = "MYAPP_"

config.Add("my-opt").
	EnvVar("MY_OPT").
	Default("default value").
	Description("The my-opt flag can be set in the MYAPP_MY_OPT env variable")

config.Parse()

```

## Disclaimer

I am new to Go, the best way to learn a language is to create a
project and expose it to the world for feedback. For a full-featured
configuration utility written by a more competent Gopher, check out
the [globalconf](https://github.com/rakyll/globalconf) project.

