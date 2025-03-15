package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPathConfigNotExists asserts that `Path` identifies non-existent config
// files
func TestPathConfigNotExists(t *testing.T) {

	// package (invalid) cheatpaths
	paths := []string{"/cheat-test-conf-does-not-exist"}
	_, err := Path(paths)
	assert.Error(t, err, "failed to identify non-existent config file")
}

// TestPathConfigExists asserts that `Path` identifies existent config files
func TestPathConfigExists(t *testing.T) {

	// initialize a temporary config file
	confFile, err := os.CreateTemp("", "cheat-test")
	assert.NoError(t, err, "failed to create temp file")

	// clean up the temp file
	defer os.Remove(confFile.Name())

	// package cheatpaths
	paths := []string{
		"/cheat-test-conf-does-not-exist",
		confFile.Name(),
	}

	// assert
	got, err := Path(paths)
	assert.NoError(t, err, "failed to identify config file")
	assert.Equal(t, confFile.Name(), got, "failed to return config path")
}
