package cmd

import (
	"fmt"
	"os"

	"github.com/yagoyudi/cheat/internal/config"
	"github.com/yagoyudi/cheat/internal/display"
	"github.com/yagoyudi/cheat/internal/sheets"
)

// cmdTags lists all tags in use.
func Tags(_ map[string]interface{}, conf config.Config) {

	// load the cheatsheets
	cheatsheets, err := sheets.Load(conf.Cheatpaths)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list cheatsheets: %v\n", err)
		os.Exit(1)
	}

	// assemble the output
	out := ""
	for _, tag := range sheets.Tags(cheatsheets) {
		out += fmt.Sprintln(tag)
	}

	// display the output
	display.Write(out, conf)
}
