package display

import (
	"fmt"
	"strings"
)

// Indent prepends each line of a string with a tab
func Indent(str string) string {
	str = strings.TrimSpace(str)

	out := ""
	for _, line := range strings.Split(str, "\n") {
		out += fmt.Sprintf("\t%s\n", line)
	}

	return out
}
