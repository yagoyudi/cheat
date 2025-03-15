package cheatpath

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFilterSuccess asserts that the proper cheatpath is returned when the
// requested cheatpath exists
func TestFilterSuccess(t *testing.T) {

	// init cheatpaths
	paths := []Cheatpath{
		{Name: "foo"},
		{Name: "bar"},
		{Name: "baz"},
	}

	// filter the paths
	paths, err := Filter(paths, "bar")
	assert.NoError(t, err)

	// assert that the expected path was returned
	assert.Equal(t, 1, len(paths))
	assert.Equal(t, "bar", paths[0].Name)
}

// TestFilterFailure asserts that an error is returned when a non-existent
// cheatpath is requested
func TestFilterFailure(t *testing.T) {

	// init cheatpaths
	paths := []Cheatpath{
		{Name: "foo"},
		{Name: "bar"},
		{Name: "baz"},
	}

	// filter the paths
	_, err := Filter(paths, "qux")
	assert.Error(t, err)
}
