// Package display implement functions pertaining to writing formatted
// cheatsheet content to stdout, or alternatively the system pager.
package display

import (
	"fmt"

	"github.com/yagoyudi/cheat/internal/config"
)

// Faint returns a faintly-colored string that's used to de-prioritize text
// written to stdout
func Faint(str string, conf config.Config) string {
	if conf.Colorize {
		return fmt.Sprintf("\033[2m%s\033[0m", str)
	}
	return str
}
