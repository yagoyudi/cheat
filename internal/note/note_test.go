package note

import (
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yagoyudi/note/internal/config"
	"github.com/yagoyudi/note/internal/mock"
)

// Asserts that notes initialize properly
func TestNoteSuccess(t *testing.T) {
	note, err := New(
		"foo",
		"community",
		mock.Path("note/foo"),
		[]string{"alpha", "bravo"},
		false,
	)
	assert.NoError(t, err, "failed to load note")
	assert.Equal(t, "foo", note.Name, "failed to init title")
	assert.Equal(t, mock.Path("note/foo"), note.Path, "failed to init path")
	wantText := "# To foo the bar:\n  foo bar\n"
	assert.Equal(t, wantText, note.Body, "failed to init text")

	// Tags should sort alphabetically:
	wantTags := []string{"alpha", "bar", "baz", "bravo", "foo"}
	assert.Equal(t, wantTags, note.Tags, "failed to init tags")
	assert.Equal(t, "sh", note.Syntax, "failed to init syntax")
	assert.Equal(t, false, note.ReadOnly, "failed to init readonly")
}

// Asserts that an error is returned if the note cannot be read
func TestNoteFailure(t *testing.T) {
	_, err := New(
		"foo",
		"community",
		mock.Path("/does-not-exist"),
		[]string{"alpha", "bravo"},
		false,
	)
	assert.Error(t, err, "failed to return an error on unreadable sheet")
}

// Asserts that an error is returned if the notes's header cannot be parsed.
func TestNoteHeaderFailure(t *testing.T) {
	_, err := New(
		"foo",
		"community",
		mock.Path("sheet/bad-header"),
		[]string{"alpha", "bravo"},
		false,
	)
	assert.Error(t, err, "failed to return an error on malformed front-matter")
}

// Asserts that raw note is properly parsed when it contains header
func TestHasHeader(t *testing.T) {
	raw := `---
syntax: go
tags: [ test ]
---
To foo the bar: baz`

	header, body, err := parse(raw)
	assert.NoError(t, err, "failed to parse markdown")

	want := "To foo the bar: baz"
	assert.Equal(t, want, body, "failed to parse text")

	want = "go"
	assert.Equal(t, want, header.Syntax, "failed to parse syntax")

	want = "test"
	assert.Equal(t, want, header.Tags[0], "failed to parse tags")
	assert.Equal(t, 1, len(header.Tags), "failed to parse tags")
}

// Asserts that way note is properly parsed when it does not contain header
func TestHasNoHeader(t *testing.T) {
	raw := "To foo the bar: baz"
	header, body, err := parse(raw)
	assert.NoError(t, err, "failed to parse raw note")
	assert.Equal(t, raw, body, "failed to parse body")
	assert.Equal(t, "", header.Syntax, "failex to parse syntax")
	assert.Equal(t, 0, len(header.Tags), "failed to parse tags")
}

// Asserts that raw note is properly parsed when it contains invalid header
func TestHasInvalidHeader(t *testing.T) {
	raw := `---
syntax: go
tags: [ test ]
To foo the bar: baz`

	_, _, err := parse(raw)
	assert.Error(t, err, "failed to error on invalid header")
}

// Ensures that the expected output is returned when no matches are found
func TestSearchNoMatch(t *testing.T) {
	note := Note{
		Body: "The quick brown fox\njumped over\nthe lazy dog.",
	}
	reg, err := regexp.Compile("(?i)foo")
	assert.NoError(t, err, "failed to compile regex")

	matches := note.Search(reg)
	assert.Equal(t, "", matches, "failure: expected no matches")
}

// Asserts that the expected output is returned when a single match is returned
func TestSearchSingleMatch(t *testing.T) {
	note := Note{
		Body: "The quick brown fox\njumped over\n\nthe lazy dog.",
	}
	reg, err := regexp.Compile("(?i)fox")
	assert.NoError(t, err, "failed to compile regex")

	matches := note.Search(reg)
	want := "The quick brown fox\njumped over"
	assert.Equal(t, want, matches, "failed to return expected matches")
}

// Asserts that the expected output is returned when a multiple matches are
// returned
func TestSearchMultiMatch(t *testing.T) {
	note := Note{
		Body: "The quick brown fox\n\njumped over\n\nthe lazy dog.",
	}

	reg, err := regexp.Compile("(?i)the")
	assert.NoError(t, err, "failed to compile regex")

	matches := note.Search(reg)
	want := "The quick brown fox\n\nthe lazy dog."
	assert.Equal(t, want, matches, "failed to return expected matches")
}

// Asserts that Copy correctly copies files at a single level of depth
func TestCopyFlat(t *testing.T) {
	body := "this is the cheatsheet text"
	src, err := os.CreateTemp("", "foo-src")
	assert.NoError(t, err, "failed to mock note")
	defer src.Close()
	defer os.Remove(src.Name())

	_, err = src.WriteString(body)
	assert.NoError(t, err, "failed to write to mock note")

	note, err := New("foo", "community", src.Name(), []string{}, false)
	assert.NoError(t, err, "failed to init cheatsheet")

	// compute the outfile's path
	outpath := path.Join(os.TempDir(), note.Name)
	defer os.Remove(outpath)

	// attempt to copy the cheatsheet
	err = note.Copy(outpath)
	assert.NoError(t, err, "failed to copy cheatsheet")

	// assert that the destination file contains the correct text
	got, err := os.ReadFile(outpath)
	assert.NoError(t, err, "failed to read destination file")
	assert.Equal(t, body, string(got), "destination file contained wrong text")
}

// Asserts that Copy correctly copies files at several levels of depth
func TestCopyDeep(t *testing.T) {
	text := "this is the cheatsheet text"
	src, err := os.CreateTemp("", "foo-src")
	assert.NoError(t, err, "failed to mock cheatsheet")
	defer src.Close()
	defer os.Remove(src.Name())

	_, err = src.WriteString(text)
	assert.NoError(t, err, "failed to write to mock cheatsheet")

	note, err := New(
		"/cheat-tests/alpha/bravo/foo",
		"community",
		src.Name(),
		[]string{},
		false,
	)
	assert.NoError(t, err, "failed to init cheatsheet")

	// compute the outfile's path
	outpath := path.Join(os.TempDir(), note.Name)
	defer os.RemoveAll(path.Join(os.TempDir(), "cheat-tests"))

	// attempt to copy the cheatsheet
	err = note.Copy(outpath)
	assert.NoError(t, err, "failed to copy cheatsheet")

	// assert that the destination file contains the correct text
	got, err := os.ReadFile(outpath)
	assert.NoError(t, err, "failed to read destination file")
	assert.Equal(t, text, string(got), "destination file contained wrong text")
}

// Asserts that syntax-highlighting is correctly applied
func TestColorize(t *testing.T) {
	conf := config.Config{
		Formatter: "terminal16m",
		Style:     "solarized-dark",
	}
	note := Note{
		Body: "echo 'foo'",
	}

	note.Colorize(conf)
	want := "[38;2;181;137;0mecho[0m[38;2;147;161;161m"
	want += " [0m[38;2;42;161;152m'foo'[0m"
	assert.Equal(t, want, note.Body, "failed to colorize sheet")
}

// Ensures that tags are properly recognized as being absent or present
func TestTagged(t *testing.T) {
	tags := []string{"foo", "bar", "baz"}
	note := Note{Tags: tags}

	// Assert that set tags are recognized as set:
	for _, tag := range tags {
		assert.NotEqual(t, false, note.Tagged(tag), "failed to recognize tag")
	}

	// Assert that unset tags are recognized as unset:
	assert.Equal(t, false, note.Tagged("qux"), "failed to recognize absent tag")
}
