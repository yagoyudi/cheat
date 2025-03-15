package sheet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHasFrontmatter asserts that markdown is properly parsed when it contains
// frontmatter
func TestHasFrontmatter(t *testing.T) {

	// stub our cheatsheet content
	markdown := `---
syntax: go
tags: [ test ]
---
To foo the bar: baz`

	// parse the frontmatter
	fm, text, err := parse(markdown)
	assert.NoError(t, err, "failed to parse markdown")

	want := "To foo the bar: baz"
	assert.Equal(t, want, text, "failed to parse text")

	want = "go"
	assert.Equal(t, want, fm.Syntax, "failed to parse syntax")

	want = "test"
	assert.Equal(t, want, fm.Tags[0], "failed to parse tags")
	assert.Equal(t, 1, len(fm.Tags), "failed to parse tags")
}

// TestHasFrontmatter asserts that markdown is properly parsed when it does not
// contain frontmatter
func TestHasNoFrontmatter(t *testing.T) {

	// stub our cheatsheet content
	markdown := "To foo the bar: baz"

	// parse the frontmatter
	fm, text, err := parse(markdown)
	assert.NoError(t, err, "failed to parse markdown")
	assert.Equal(t, markdown, text, "failed to parse text")
	assert.Equal(t, "", fm.Syntax, "failex to parse syntax")
	assert.Equal(t, 0, len(fm.Tags), "failed to parse tags")
}

// TestHasInvalidFrontmatter asserts that markdown is properly parsed when it
// contains invalid frontmatter
func TestHasInvalidFrontmatter(t *testing.T) {

	// stub our cheatsheet content (with invalid frontmatter)
	markdown := `---
syntax: go
tags: [ test ]
To foo the bar: baz`

	// parse the frontmatter
	_, text, err := parse(markdown)

	// assert that an error was returned
	assert.Error(t, err, "failed to error on invalid frontmatter")
	assert.Equal(t, markdown, text, "failed to parse text")
}
