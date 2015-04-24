package gofigure

import (
	"os"
	"testing"
)

const envPrefixTest = "TEST_"
const jsonTestFixture = "test_fixture.json"
const fileSpecTest = "main_category.listen"

func TestEnv(t *testing.T) {
	os.Setenv(envPrefixTest+"LISTEN", ":3001")

	config := New()
	config.EnvPrefix = envPrefixTest

	config.Add("listen").EnvVar("LISTEN").Default(":3000")

	config.DisableCommandLine = true
	config.Parse()

	listen := *config.Get("listen")
	if listen != ":3001" {
		t.Errorf("Flag 'listen' found %v, expected :3001", listen)
	}
}

func TestJsonFile(t *testing.T) {
	config := New()
	config.FileParser = NewJsonFile(jsonTestFixture)

	config.Add("listen").FileSpec(fileSpecTest).Default(":3000")
	config.DisableCommandLine = true
	config.Parse()

	listen := *config.Get("listen")
	if listen != ":3002" {
		t.Errorf("Flag 'listen' found %v, expected :3002", listen)
	}
}

func TestPrecedenceEnv(t *testing.T) {
	os.Setenv(envPrefixTest+"LISTEN", ":3001")
	config := New()
	config.FileParser = NewJsonFile(jsonTestFixture)
	config.EnvOverridesFile = true
	config.EnvPrefix = envPrefixTest

	config.Add("listen").FileSpec(fileSpecTest).EnvVar("LISTEN").Default(":3000")
	config.DisableCommandLine = true
	config.Parse()

	listen := *config.Get("listen")
	if listen != ":3001" {
		t.Errorf("Flag 'listen' found %v, expected :3001", listen)
	}
}

func TestPrecedenceFile(t *testing.T) {
	os.Setenv(envPrefixTest+"LISTEN", ":3001")
	config := New()
	config.FileParser = NewJsonFile(jsonTestFixture)
	config.EnvOverridesFile = false
	config.EnvPrefix = envPrefixTest

	config.Add("listen").FileSpec(fileSpecTest).EnvVar("LISTEN").Default(":3000")
	config.DisableCommandLine = true
	config.Parse()

	listen := *config.Get("listen")
	if listen != ":3002" {
		t.Errorf("Flag 'listen' found %v, expected :3002", listen)
	}
}

