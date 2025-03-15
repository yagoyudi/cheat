package cheatpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestWriteableOK asserts that Writeable returns the appropriate cheatpath
// when a writeable cheatpath exists
func TestWriteableOK(t *testing.T) {

	// initialize some cheatpaths
	cheatpaths := []Cheatpath{
		{Path: "/foo", ReadOnly: true},
		{Path: "/bar", ReadOnly: false},
		{Path: "/baz", ReadOnly: true},
	}

	// get the writeable cheatpath
	got, err := Writeable(cheatpaths)
	assert.NoError(t, err)
	assert.Equal(t, "/bar", got.Path)
	assert.Equal(t, false, got.ReadOnly)
}

// TestWriteableOK asserts that Writeable returns an error when no writeable
// cheatpaths exist
func TestWriteableNotOK(t *testing.T) {

	// initialize some cheatpaths
	cheatpaths := []Cheatpath{
		{Path: "/foo", ReadOnly: true},
		{Path: "/bar", ReadOnly: true},
		{Path: "/baz", ReadOnly: true},
	}

	// get the writeable cheatpath
	_, err := Writeable(cheatpaths)
	assert.Error(t, err, "failed to return an error when no writeable paths found")
}
