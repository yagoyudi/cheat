package sheets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yagoyudi/cheat/internal/sheet"
)

// TestSort asserts that Sort properly sorts sheets
func TestSort(t *testing.T) {

	// mock a map of cheatsheets
	sheets := map[string]sheet.Sheet{
		"foo": {Title: "foo"},
		"bar": {Title: "bar"},
		"baz": {Title: "baz"},
	}

	// sort the sheets
	sorted := Sort(sheets)

	// assert that the sheets sorted properly
	want := []string{"bar", "baz", "foo"}
	for i, got := range sorted {
		assert.Equal(t, want[i], got.Title, "sort returned incorrect value")
	}
}
