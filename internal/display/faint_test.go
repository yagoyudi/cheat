package display

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yagoyudi/cheat/internal/config"
)

// TestFaint asserts that Faint applies faint formatting
func TestFaint(t *testing.T) {

	// case: apply colorization
	conf := config.Config{Colorize: true}
	want := "\033[2mfoo\033[0m"
	got := Faint("foo", conf)
	assert.Equal(t, want, got, "failed to faint")

	// case: do not apply colorization
	conf.Colorize = false
	want = "foo"
	got = Faint("foo", conf)
	assert.Equal(t, want, got, "failed to faint")
}
