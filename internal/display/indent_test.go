package display

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIndent asserts that Indent prepends a tab to each line
func TestIndent(t *testing.T) {
	got := Indent("foo\nbar\nbaz")
	want := "\tfoo\n\tbar\n\tbaz\n"
	assert.Equal(t, want, got, "failed to indent")
}
