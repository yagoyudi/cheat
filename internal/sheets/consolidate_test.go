package sheets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yagoyudi/cheat/internal/sheet"
)

// TestConsolidate asserts that cheatsheets are properly consolidated
func TestConsolidate(t *testing.T) {

	// mock cheatsheets available on multiple cheatpaths
	cheatpaths := []map[string]sheet.Sheet{

		// mock community cheatsheets
		{
			"foo": {Title: "foo", Path: "community/foo"},
			"bar": {Title: "bar", Path: "community/bar"},
		},

		// mock local cheatsheets
		{
			"bar": {Title: "bar", Path: "local/bar"},
			"baz": {Title: "baz", Path: "local/baz"},
		},
	}

	// consolidate the cheatsheets
	consolidated := Consolidate(cheatpaths)

	// specify the expected output
	want := map[string]sheet.Sheet{
		"foo": {Title: "foo", Path: "community/foo"},
		"bar": {Title: "bar", Path: "local/bar"},
		"baz": {Title: "baz", Path: "local/baz"},
	}
	assert.Equal(t, want, consolidated, "failed to consolidate cheatpaths")
}
