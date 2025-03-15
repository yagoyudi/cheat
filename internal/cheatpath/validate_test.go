package cheatpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateValid asserts that valid cheatpaths validate successfully
func TestValidateValid(t *testing.T) {

	// initialize a valid cheatpath
	cheatpath := Cheatpath{
		Name:     "foo",
		Path:     "/foo",
		ReadOnly: false,
		Tags:     []string{},
	}

	err := cheatpath.Validate()
	assert.NoError(t, err)
}

// TestValidateMissingName asserts that paths that are missing a name fail to
// validate
func TestValidateMissingName(t *testing.T) {

	// initialize a valid cheatpath
	cheatpath := Cheatpath{
		Path:     "/foo",
		ReadOnly: false,
		Tags:     []string{},
	}

	err := cheatpath.Validate()
	assert.Error(t, err)
}

// TestValidateMissingPath asserts that paths that are missing a path fail to
// validate
func TestValidateMissingPath(t *testing.T) {

	// initialize a valid cheatpath
	cheatpath := Cheatpath{
		Name:     "foo",
		ReadOnly: false,
		Tags:     []string{},
	}

	err := cheatpath.Validate()
	assert.Error(t, err)
}
