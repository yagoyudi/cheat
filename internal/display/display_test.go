package display

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yagoyudi/cheat/internal/config"
)

// Asserts that Faint applies faint formatting
func TestFaint(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		config   config.Config
		expected string
	}{
		{
			name:     "apply colorization",
			input:    "foo",
			config:   config.Config{Colorize: true},
			expected: "\033[2mfoo\033[0m",
		},
		{
			name:     "do not apply colorization",
			input:    "foo",
			config:   config.Config{Colorize: false},
			expected: "foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Faint(tt.input, tt.config)
			assert.Equal(t, got, tt.expected)
		})
	}
}

// Asserts that Indent prepends a tab to each line
func TestIndent(t *testing.T) {
	got := Indent("foo\nbar\nbaz")
	want := "\tfoo\n\tbar\n\tbaz\n"
	assert.Equal(t, want, got, "failed to indent")
}
