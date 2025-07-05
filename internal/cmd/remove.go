package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/note/internal/config"
	"github.com/yagoyudi/note/internal/notes"
)

func init() {
	removeCmd.Flags().StringP("tag", "t", "", "filter notes by tag")
}

var removeCmd = &cobra.Command{
	Use:     "rm [note]",
	Short:   "Removes a note",
	Args:    cobra.ExactArgs(1),
	Example: `  note rm kubectl -t community`,
	Run: func(cmd *cobra.Command, args []string) {
		tags, err := cmd.Flags().GetString("tag")
		cobra.CheckErr(err)

		var conf config.Config
		cobra.CheckErr(viper.Unmarshal(&conf))

		loadedNotes, err := notes.Load(conf.Notebooks)
		cobra.CheckErr(err)

		if cmd.Flags().Changed("tag") {
			loadedNotes = notes.Filter(loadedNotes, strings.Split(tags, ","))
		}

		noteName := args[0]
		consolidatedNotes := notes.Consolidate(loadedNotes)
		note, ok := consolidatedNotes[noteName]
		if !ok {
			fmt.Printf("No cheatsheet found for '%s'.\n", noteName)
			os.Exit(0)
		}

		if note.ReadOnly {
			fmt.Printf("Cheatsheet '%s' is read-only.\n", noteName)
			os.Exit(0)
		}

		cobra.CheckErr(os.Remove(note.Path))
	},
}
