package notebook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Asserts that valid notebook validate successfully
func TestValidateValid(t *testing.T) {
	notebook := Notebook{
		Name:     "foo",
		Path:     "/foo",
		ReadOnly: false,
		Tags:     []string{},
	}
	err := notebook.Validate()
	assert.NoError(t, err)
}

// Asserts that notebooks that are missing a name fail to validate
func TestValidateMissingName(t *testing.T) {
	notebook := Notebook{
		Path:     "/foo",
		ReadOnly: false,
		Tags:     []string{},
	}
	err := notebook.Validate()
	assert.Error(t, err)
}

// Asserts that notebooks that are missing a path fail to validate
func TestValidateMissingPath(t *testing.T) {
	notebook := Notebook{
		Name:     "foo",
		ReadOnly: false,
		Tags:     []string{},
	}
	err := notebook.Validate()
	assert.Error(t, err)
}

// Asserts that the proper notebook is returned when the requested notebook
// exists
func TestFilterSuccess(t *testing.T) {
	notebooks := []Notebook{
		{Name: "foo"},
		{Name: "bar"},
		{Name: "baz"},
	}

	notebooks, err := Filter(notebooks, "bar")
	assert.NoError(t, err)

	assert.Equal(t, 1, len(notebooks))
	assert.Equal(t, "bar", notebooks[0].Name)
}

// Asserts that an error is returned when a non-existent notebook is requested
func TestFilterFailure(t *testing.T) {
	notebooks := []Notebook{
		{Name: "foo"},
		{Name: "bar"},
		{Name: "baz"},
	}
	_, err := Filter(notebooks, "qux")
	assert.Error(t, err)
}

// Asserts that Writeable returns the appropriate notebook when a writeable
// notebook exists
func TestWriteableOK(t *testing.T) {
	notebooks := []Notebook{
		{Path: "/foo", ReadOnly: true},
		{Path: "/bar", ReadOnly: false},
		{Path: "/baz", ReadOnly: true},
	}
	got, err := Writeable(notebooks)
	assert.NoError(t, err)
	assert.Equal(t, "/bar", got.Path)
	assert.Equal(t, false, got.ReadOnly)
}

// Asserts that Writeable returns an error when no writeable notebook exist
func TestWriteableNotOK(t *testing.T) {
	notebooks := []Notebook{
		{Path: "/foo", ReadOnly: true},
		{Path: "/bar", ReadOnly: true},
		{Path: "/baz", ReadOnly: true},
	}
	_, err := Writeable(notebooks)
	assert.Error(t, err, "failed to return an error when no writeable notebook found")
}
