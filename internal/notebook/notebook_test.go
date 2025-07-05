package notebook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Asserts that valid notepaths validate successfully
func TestValidateValid(t *testing.T) {
	notepath := Notebook{
		Name:     "foo",
		Path:     "/foo",
		ReadOnly: false,
		Tags:     []string{},
	}
	err := notepath.Validate()
	assert.NoError(t, err)
}

// Asserts that paths that are missing a name fail to validate
func TestValidateMissingName(t *testing.T) {
	notepath := Notebook{
		Path:     "/foo",
		ReadOnly: false,
		Tags:     []string{},
	}
	err := notepath.Validate()
	assert.Error(t, err)
}

// Asserts that paths that are missing a path fail to validate
func TestValidateMissingPath(t *testing.T) {
	notepath := Notebook{
		Name:     "foo",
		ReadOnly: false,
		Tags:     []string{},
	}
	err := notepath.Validate()
	assert.Error(t, err)
}

// Asserts that the proper notepath is returned when the requested notepath
// exists
func TestFilterSuccess(t *testing.T) {
	paths := []Notebook{
		{Name: "foo"},
		{Name: "bar"},
		{Name: "baz"},
	}

	paths, err := Filter(paths, "bar")
	assert.NoError(t, err)

	assert.Equal(t, 1, len(paths))
	assert.Equal(t, "bar", paths[0].Name)
}

// Asserts that an error is returned when a non-existent notepath is requested
func TestFilterFailure(t *testing.T) {
	paths := []Notebook{
		{Name: "foo"},
		{Name: "bar"},
		{Name: "baz"},
	}
	_, err := Filter(paths, "qux")
	assert.Error(t, err)
}

// Asserts that Writeable returns the appropriate notepath when a writeable
// notepath exists
func TestWriteableOK(t *testing.T) {
	notepaths := []Notebook{
		{Path: "/foo", ReadOnly: true},
		{Path: "/bar", ReadOnly: false},
		{Path: "/baz", ReadOnly: true},
	}
	got, err := Writeable(notepaths)
	assert.NoError(t, err)
	assert.Equal(t, "/bar", got.Path)
	assert.Equal(t, false, got.ReadOnly)
}

// Asserts that Writeable returns an error when no writeable notepaths exist
func TestWriteableNotOK(t *testing.T) {
	notepaths := []Notebook{
		{Path: "/foo", ReadOnly: true},
		{Path: "/bar", ReadOnly: true},
		{Path: "/baz", ReadOnly: true},
	}
	_, err := Writeable(notepaths)
	assert.Error(t, err, "failed to return an error when no writeable paths found")
}
