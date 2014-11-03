# gofigure

A configuration utility for Go inspired by [The Twelve-Factor App](http://12factor.net/config) methodology.


Configuration options are defined via command line arguments and
environment variables, evaluated in that order.

## Usage

```go
import "github.com/cpliakas/gofigure"
```

```go
config := gofigure.New()
config.EnvPrefix = "MYAPP_"

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

