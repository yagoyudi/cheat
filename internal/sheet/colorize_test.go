package sheet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yagoyudi/cheat/internal/config"
)

// TestColorize asserts that syntax-highlighting is correctly applied
func TestColorize(t *testing.T) {

	// mock configs
	conf := config.Config{
		Formatter: "terminal16m",
		Style:     "solarized-dark",
	}

	// mock a sheet
	s := Sheet{
		Text: "echo 'foo'",
	}

	// colorize the sheet text
	s.Colorize(conf)

	// initialize expectations
	want := "[38;2;181;137;0mecho[0m[38;2;147;161;161m"
	want += " [0m[38;2;42;161;152m'foo'[0m"

	// assert
	assert.Equal(t, want, s.Text, "failed to colorize sheet")
}
