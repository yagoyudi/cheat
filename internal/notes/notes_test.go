package notes

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yagoyudi/note/internal/mock"
	"github.com/yagoyudi/note/internal/note"
	"github.com/yagoyudi/note/internal/notebook"
)

// Asserts that notes are properly consolidated
func TestConsolidate(t *testing.T) {
	notes := []map[string]note.Note{
		{
			"foo": {Name: "foo", Path: "community/foo"},
			"bar": {Name: "bar", Path: "community/bar"},
		},
		{
			"bar": {Name: "bar", Path: "local/bar"},
			"baz": {Name: "baz", Path: "local/baz"},
		},
	}

	consolidated := Consolidate(notes)
	want := map[string]note.Note{
		"foo": {Name: "foo", Path: "community/foo"},
		"bar": {Name: "bar", Path: "local/bar"},
		"baz": {Name: "baz", Path: "local/baz"},
	}
	assert.Equal(t, want, consolidated, "failed to consolidate cheatpaths")
}

// Asserts that Sort properly sorts notes
func TestSort(t *testing.T) {
	notes := map[string]note.Note{
		"foo": {Name: "foo"},
		"bar": {Name: "bar"},
		"baz": {Name: "baz"},
	}
	sorted := Sort(notes)
	want := []string{"bar", "baz", "foo"}
	for i, got := range sorted {
		assert.Equal(t, want[i], got.Name, "sort returned incorrect value")
	}
}

// Asserts that note tags are properly returned
func TestTags(t *testing.T) {
	notepaths := []map[string]note.Note{
		{
			"foo": {Name: "foo", Tags: []string{"alpha"}},
			"bar": {Name: "bar", Tags: []string{"alpha", "bravo"}},
		},
		{
			"bar": {Name: "bar", Tags: []string{"bravo", "charlie"}},
			"baz": {Name: "baz", Tags: []string{"delta"}},
		},
	}
	tags := Tags(notepaths)
	want := []string{
		"alpha",
		"bravo",
		"charlie",
		"delta",
	}
	assert.Equal(t, want, tags, "failed to return tags")
}

// Asserts that Filter properly filters results when passed a single tag
func TestFilterSingleTag(t *testing.T) {
	notepaths := []map[string]note.Note{
		{
			"foo": {Name: "foo", Tags: []string{"alpha", "bravo"}},
			"bar": {Name: "bar", Tags: []string{"bravo", "charlie"}},
		},
		{
			"baz": {Name: "baz", Tags: []string{"alpha", "bravo"}},
			"bat": {Name: "bat", Tags: []string{"bravo", "charlie"}},
		},
	}
	filtered := Filter(notepaths, []string{"bravo"})
	want := []map[string]note.Note{
		{
			"foo": {Name: "foo", Tags: []string{"alpha", "bravo"}},
			"bar": {Name: "bar", Tags: []string{"bravo", "charlie"}},
		},

		{
			"baz": {Name: "baz", Tags: []string{"alpha", "bravo"}},
			"bat": {Name: "bat", Tags: []string{"bravo", "charlie"}},
		},
	}

	assert.Equal(t, want, filtered, "failed to return expected results")
}

// Asserts that Filter properly filters results when passed multiple tags
func TestFilterMultiTag(t *testing.T) {
	notepaths := []map[string]note.Note{
		{
			"foo": {Name: "foo", Tags: []string{"alpha", "bravo"}},
			"bar": {Name: "bar", Tags: []string{"bravo", "charlie"}},
		},

		{
			"baz": {Name: "baz", Tags: []string{"alpha", "bravo"}},
			"bat": {Name: "bat", Tags: []string{"bravo", "charlie"}},
		},
	}
	filtered := Filter(notepaths, []string{"alpha", "bravo"})
	want := []map[string]note.Note{
		{
			"foo": {Name: "foo", Tags: []string{"alpha", "bravo"}},
		},
		{
			"baz": {Name: "baz", Tags: []string{"alpha", "bravo"}},
		},
	}
	assert.Equal(t, want, filtered, "failed to return expected result")
}

// Asserts that notes on valid notepaths can be loaded successfully
func TestLoad(t *testing.T) {
	notepaths := []notebook.Notebook{
		{
			Name:     "community",
			Path:     path.Join(mock.Path("cheatsheets"), "community"),
			ReadOnly: true,
		},
		{
			Name:     "personal",
			Path:     path.Join(mock.Path("cheatsheets"), "personal"),
			ReadOnly: false,
		},
	}

	notes, err := Load(notepaths)
	assert.NoError(t, err, "failed to load cheatsheets")

	// Assert that the correct number of note loaded:
	want := 4
	assert.Equal(t, want, len(notes), "failed to load correct number of cheatsheets")
}

// Asserts that an error is returned if a path is invalid
func TestLoadBadPath(t *testing.T) {
	notepaths := []notebook.Notebook{
		{
			Name:     "badpath",
			Path:     "/cheat/test/path/does/not/exist",
			ReadOnly: true,
		},
	}
	_, err := Load(notepaths)
	assert.Error(t, err, "failed to reject invalid cheatpath")
}
