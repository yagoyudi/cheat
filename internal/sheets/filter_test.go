package sheets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yagoyudi/cheat/internal/sheet"
)

// TestFilterSingleTag asserts that Filter properly filters results when passed
// a single tag
func TestFilterSingleTag(t *testing.T) {

	// mock cheatsheets available on multiple cheatpaths
	cheatpaths := []map[string]sheet.Sheet{

		{
			"foo": {Title: "foo", Tags: []string{"alpha", "bravo"}},
			"bar": {Title: "bar", Tags: []string{"bravo", "charlie"}},
		},

		{
			"baz": {Title: "baz", Tags: []string{"alpha", "bravo"}},
			"bat": {Title: "bat", Tags: []string{"bravo", "charlie"}},
		},
	}

	// filter the cheatsheets
	filtered := Filter(cheatpaths, []string{"bravo"})

	// assert that the expect results were returned
	want := []map[string]sheet.Sheet{
		{
			"foo": {Title: "foo", Tags: []string{"alpha", "bravo"}},
			"bar": {Title: "bar", Tags: []string{"bravo", "charlie"}},
		},

		{
			"baz": {Title: "baz", Tags: []string{"alpha", "bravo"}},
			"bat": {Title: "bat", Tags: []string{"bravo", "charlie"}},
		},
	}

	assert.Equal(t, want, filtered, "failed to return expected results")
}

// TestFilterSingleTag asserts that Filter properly filters results when passed
// multiple tags
func TestFilterMultiTag(t *testing.T) {

	// mock cheatsheets available on multiple cheatpaths
	cheatpaths := []map[string]sheet.Sheet{

		{
			"foo": {Title: "foo", Tags: []string{"alpha", "bravo"}},
			"bar": {Title: "bar", Tags: []string{"bravo", "charlie"}},
		},

		{
			"baz": {Title: "baz", Tags: []string{"alpha", "bravo"}},
			"bat": {Title: "bat", Tags: []string{"bravo", "charlie"}},
		},
	}

	// filter the cheatsheets
	filtered := Filter(cheatpaths, []string{"alpha", "bravo"})

	// assert that the expect results were returned
	want := []map[string]sheet.Sheet{
		{
			"foo": {Title: "foo", Tags: []string{"alpha", "bravo"}},
		},
		{
			"baz": {Title: "baz", Tags: []string{"alpha", "bravo"}},
		},
	}

	assert.Equal(t, want, filtered, "failed to return expected result")
}
