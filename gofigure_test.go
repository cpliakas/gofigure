package gofigure

import (
	"os"
	"testing"
)

const envPrefixTest = "TEST_"

func TestParse(t *testing.T) {
	os.Setenv(envPrefixTest+"LISTEN", ":3001")

	config := New()
	config.EnvPrefix = envPrefixTest

	config.Add("listen").EnvVar("LISTEN").Default(":3000")

	config.DisableCommandLine = true
	config.Parse()

	listen := config.Get("listen")
	if listen != ":3001" {
		t.Errorf("Flag 'listen' found %v, expected :3001", listen)
	}
}
