package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestColor asserts that colorization rules are properly respected
func TestColor(t *testing.T) {

	noColorConf := Config{
		Colorize: false,
	}
	assert.Equal(t, false, noColorConf.Colorize)

	colorConf := Config{
		Colorize: true,
	}
	assert.Equal(t, true, colorConf.Colorize)
}
