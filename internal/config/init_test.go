package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInit asserts that configs are properly initialized
func TestInit(t *testing.T) {

	// initialize a temporary config file
	confFile, err := os.CreateTemp("", "cheat-test")
	assert.NoError(t, err, "failed to create temp file")

	// clean up the temp file
	defer os.Remove(confFile.Name())

	// initialize the config file
	conf := "mock config data"
	err = Init(confFile.Name(), conf)
	assert.NoError(t, err, "failed to init config file")

	// read back the config file contents
	bytes, err := os.ReadFile(confFile.Name())
	assert.NoError(t, err, "failed to read config file")

	// assert that the contents were written correctly
	got := string(bytes)
	assert.Equal(t, conf, got, "failed to write config")
}
