package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/cheat/internal/config"
	"github.com/yagoyudi/cheat/internal/display"
	"github.com/yagoyudi/cheat/internal/sheets"
)

var tagsCmd = &cobra.Command{
	Use:     `tags`,
	Short:   `Lists all tags in use`,
	Example: `  cheat tags`,
	RunE: func(_ *cobra.Command, _ []string) error {
		var conf config.Config
		if err := viper.Unmarshal(&conf); err != nil {
			return err
		}

		// load the cheatsheets
		cheatsheets, err := sheets.Load(conf.Cheatpaths)
		if err != nil {
			return err
		}

		// assemble the output
		out := ""
		for _, tag := range sheets.Tags(cheatsheets) {
			out += fmt.Sprintln(tag)
		}

		// display the output
		display.Write(out, conf)

		return nil
	},
}
