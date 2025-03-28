package config

import (
	"os"

	"github.com/mattn/go-isatty"
)

// Color indicates whether colorization should be applied to the output
func (c *Config) Color() bool {

	// default to the colorization specified in the configs...
	colorize := c.Colorize

	// ... however, only apply colorization if we're writing to a tty...
	if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		colorize = false
	}

	return colorize
}
