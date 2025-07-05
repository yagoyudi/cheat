package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/cheat/internal/config"
	"github.com/yagoyudi/cheat/internal/display"
	"github.com/yagoyudi/cheat/internal/notes"
)

var tagsCmd = &cobra.Command{
	Use:     `tags`,
	Aliases: []string{"t"},
	Short:   `Lists all tags in use`,
	Example: `  cheat tags`,
	Run: func(_ *cobra.Command, _ []string) {
		var conf config.Config
		cobra.CheckErr(viper.Unmarshal(&conf))

		loadedNotes, err := notes.Load(conf.Notebooks)
		cobra.CheckErr(err)

		out := ""
		for _, tag := range notes.Tags(loadedNotes) {
			out += fmt.Sprintln(tag)
		}

		display.Write(out, conf)
	},
}
