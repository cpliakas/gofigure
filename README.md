# gofigure

[![Build Status](https://travis-ci.org/cpliakas/gofigure.svg)](https://travis-ci.org/cpliakas/gofigure)
[![GoDoc](https://godoc.org/github.com/cpliakas/gofigure?status.svg)](https://godoc.org/github.com/cpliakas/gofigure)

A configuration utility for Go inspired by [The Twelve-Factor App](http://12factor.net/config)
methodology that allows configuration options to be set via command
line flags and environment variables.

## Usage

```go
import "github.com/cpliakas/gofigure"
```

```go
config := gofigure.New()
config.EnvPrefix = "MYAPP_"

config.Add("listen").
	EnvVar("LISTEN").
	Default(":3000").
	Description("The address to listen on.")

config.Parse()

fmt.Println(*config.Get("listen"))
```

The following commands set the value of the `listen` flag to `:3001`:

```
./myapp --listen=:3001
MYAPP_LISTEN=:3001 ./myapp
MYAPP_LISTEN=:3002 ./myapp --listen=:3001
```

## Disclaimer

I am new to Go, the best way to learn a language is to create a
project and expose it to the world for feedback. For a full-featured
configuration utility written by a more competent Gophers, check out
the [globalconf](https://github.com/rakyll/globalconf) and
[Viper](https://github.com/spf13/viper) projects.
