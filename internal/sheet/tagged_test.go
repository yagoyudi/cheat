package sheet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTagged ensures that tags are properly recognized as being absent or
// present
func TestTagged(t *testing.T) {

	// initialize a cheatsheet
	tags := []string{"foo", "bar", "baz"}
	sheet := Sheet{Tags: tags}

	// assert that set tags are recognized as set
	for _, tag := range tags {
		assert.NotEqual(t, false, sheet.Tagged(tag), "failed to recognize tag")
	}

	// assert that unset tags are recognized as unset
	assert.Equal(t, false, sheet.Tagged("qux"), "failed to recognize absent tag")
}
