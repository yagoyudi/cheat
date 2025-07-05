// Implement functions pertaining to writing formatted note content to stdout,
// or alternatively the system pager.
package display

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/yagoyudi/note/internal/config"
)

// Returns a faintly-colored string that's used to de-prioritize text written
// to stdout
func Faint(str string, conf config.Config) string {
	if conf.Colorize {
		return fmt.Sprintf("\033[2m%s\033[0m", str)
	}
	return str
}

// Prepends each line of a string with a tab
func Indent(str string) string {
	str = strings.TrimSpace(str)
	out := ""
	for _, line := range strings.Split(str, "\n") {
		out += fmt.Sprintf("\t%s\n", line)
	}
	return out
}

// Writes output either directly to stdout, or through a pager, depending upon
// configuration.
func Write(out string, conf config.Config) {
	// If no pager was configured, print the output to stdout and exit:
	if conf.Pager == "" {
		fmt.Print(out)
		os.Exit(0)
	}

	// Pipe output through the pager
	parts := strings.Split(conf.Pager, " ")
	pager := parts[0]
	args := parts[1:]

	cmd := exec.Command(pager, args...)
	cmd.Stdin = strings.NewReader(out)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write to pager: %v\n", err)
		os.Exit(1)
	}
}
