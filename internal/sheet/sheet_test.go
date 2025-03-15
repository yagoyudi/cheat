package sheet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yagoyudi/cheat/internal/mock"
)

// TestSheetSuccess asserts that sheets initialize properly
func TestSheetSuccess(t *testing.T) {

	sheet, err := New(
		"foo",
		"community",
		mock.Path("sheet/foo"),
		[]string{"alpha", "bravo"},
		false,
	)
	assert.NoError(t, err, "failed to load sheet")
	assert.Equal(t, "foo", sheet.Title, "failed to init title")
	assert.Equal(t, mock.Path("sheet/foo"), sheet.Path, "failed to init path")
	wantText := "# To foo the bar:\n  foo bar\n"
	assert.Equal(t, wantText, sheet.Text, "failed to init text")

	// Tags should sort alphabetically
	wantTags := []string{"alpha", "bar", "baz", "bravo", "foo"}
	assert.Equal(t, wantTags, sheet.Tags, "failed to init tags")
	assert.Equal(t, "sh", sheet.Syntax, "failed to init syntax")
	assert.Equal(t, false, sheet.ReadOnly, "failed to init readonly")
}

// TestSheetFailure asserts that an error is returned if the sheet cannot be
// read
func TestSheetFailure(t *testing.T) {

	// initialize a sheet
	_, err := New(
		"foo",
		"community",
		mock.Path("/does-not-exist"),
		[]string{"alpha", "bravo"},
		false,
	)
	assert.Error(t, err, "failed to return an error on unreadable sheet")
}

// TestSheetFrontMatterFailure asserts that an error is returned if the sheet's
// frontmatter cannot be parsed.
func TestSheetFrontMatterFailure(t *testing.T) {

	// initialize a sheet
	_, err := New(
		"foo",
		"community",
		mock.Path("sheet/bad-fm"),
		[]string{"alpha", "bravo"},
		false,
	)
	assert.Error(t, err, "failed to return an error on malformed front-matter")
}
