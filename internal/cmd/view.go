package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/note/internal/config"
	"github.com/yagoyudi/note/internal/display"
	"github.com/yagoyudi/note/internal/notes"
)

func init() {
	viewCmd.Flags().BoolP("all", "A", false, "display notes from all notebooks")
	viewCmd.Flags().StringP("tag", "t", "", "filter notes by tag")
}

var viewCmd = &cobra.Command{
	Use:     "view [note]",
	Aliases: []string{"v"},
	Short:   "Displays a note for viewing",
	Args:    cobra.ExactArgs(1),
	Example: `  note view kubectl
  note v kubectl -t community`,
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

		if cmd.Flags().Changed("tag") {
			loadedNotes = notes.Filter(loadedNotes, strings.Split(tags, ","))
		}

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

		consolidatedNotes := notes.Consolidate(loadedNotes)
		note, ok := consolidatedNotes[noteName]
		if !ok {
			fmt.Printf("No note found for '%s'\n", noteName)
			os.Exit(0)
		}
		if conf.Color() {
			note.Colorize(conf)
		}
		display.Write(note.Body, conf)
	},
}
