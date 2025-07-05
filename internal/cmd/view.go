package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/cheat/internal/config"
	"github.com/yagoyudi/cheat/internal/display"
	"github.com/yagoyudi/cheat/internal/notes"
)

func init() {
	viewCmd.Flags().BoolP("all", "A", false, "display cheatsheets from all cheatpaths")
	viewCmd.Flags().StringP("tag", "t", "", "filter cheatsheets by tag")
}

var viewCmd = &cobra.Command{
	Use:     "view [cheatsheet]",
	Aliases: []string{"v"},
	Short:   "Displays a cheatsheet for viewing",
	Args:    cobra.ExactArgs(1),
	Example: `  cheat view kubectl
  cheat view kubectl -t community`,
	Run: func(cmd *cobra.Command, args []string) {
		noteName := args[0]

		allFlag, err := cmd.Flags().GetBool("all")
		cobra.CheckErr(err)

		tags, err := cmd.Flags().GetString("tag")
		cobra.CheckErr(err)

		var conf config.Config
		cobra.CheckErr(viper.Unmarshal(&conf))

		loadedNotes, err := notes.Load(conf.Notebooks)
		cobra.CheckErr(err)

		// Filter notes by tag if --tag was provided:
		if cmd.Flags().Changed("tag") {
			loadedNotes = notes.Filter(loadedNotes, strings.Split(tags, ","))
		}

		// If --all was passed, display notes from all notepaths:
		if allFlag {
			out := ""
			for _, noteByName := range loadedNotes {
				note, ok := noteByName[noteName]
				if !ok {
					continue
				}

				out += fmt.Sprintf("%s %s\n", note.Name, display.Faint(fmt.Sprintf("(%s)", note.Notebook), conf))
				if conf.Color() {
					note.Colorize(conf)
				}
				out += display.Indent(note.Body) + "\n"
			}
			display.Write(strings.TrimSuffix(out, "\n"), conf)
			os.Exit(0)
		}

		// Consolidate the notes found on all paths into a single map of
		// `title` => `note` (ie, allow more local notes to override less
		// local notes):
		consolidatedNotes := notes.Consolidate(loadedNotes)

		// Fail early if the requested note does not exist:
		note, ok := consolidatedNotes[noteName]
		if !ok {
			fmt.Printf("Error: no cheatsheet found for '%s'\n", noteName)
			os.Exit(0)
		}
		if conf.Color() {
			note.Colorize(conf)
		}
		display.Write(note.Body, conf)
	},
}
