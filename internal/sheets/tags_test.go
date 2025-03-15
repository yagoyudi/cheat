package sheets

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yagoyudi/cheat/internal/sheet"
)

// TestTags asserts that cheatsheet tags are properly returned
func TestTags(t *testing.T) {

	// mock cheatsheets available on multiple cheatpaths
	cheatpaths := []map[string]sheet.Sheet{

		// mock community cheatsheets
		{
			"foo": {Title: "foo", Tags: []string{"alpha"}},
			"bar": {Title: "bar", Tags: []string{"alpha", "bravo"}},
		},

		// mock local cheatsheets
		{
			"bar": {Title: "bar", Tags: []string{"bravo", "charlie"}},
			"baz": {Title: "baz", Tags: []string{"delta"}},
		},
	}

	// consolidate the cheatsheets
	tags := Tags(cheatpaths)

	// specify the expected output
	want := []string{
		"alpha",
		"bravo",
		"charlie",
		"delta",
	}

	// assert that the cheatsheets properly consolidated
	assert.Equal(t, want, tags, "failed to return tags")
}
