package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/note/internal/config"
	"github.com/yagoyudi/note/internal/display"
	"github.com/yagoyudi/note/internal/notebook"
	"github.com/yagoyudi/note/internal/notes"
)

func init() {
	searchCmd.Flags().StringP("tag", "t", "", "filter notes by tag")
	searchCmd.Flags().StringP("note", "n", "", "constrain the search only to matching notes")
	searchCmd.Flags().BoolP("regex", "r", false, "treat search [phrase] as a regex")
	searchCmd.Flags().StringP("book", "b", "", "filter the notebooks")
}

var searchCmd = &cobra.Command{
	Use:     "search [phrase]",
	Aliases: []string{"s"},
	Short:   "Searches for strings in notes",
	Args:    cobra.ExactArgs(1),
	Example: `  note search '(?:[0-9]{1,3}\.){3}[0-9]{1,3}' -p personal -t networking -r`,
	Run: func(cmd *cobra.Command, args []string) {
		phrase := args[0]

		regexFlag, err := cmd.Flags().GetBool("regex")
		cobra.CheckErr(err)

		noteName, err := cmd.Flags().GetString("note")
		cobra.CheckErr(err)

		var conf config.Config
		cobra.CheckErr(viper.Unmarshal(&conf))

		if cmd.Flags().Changed("book") {
			path, err := cmd.Flags().GetString("path")
			cobra.CheckErr(err)

			conf.Notebooks, err = notebook.Filter(conf.Notebooks, path)
			cobra.CheckErr(err)
		}

		loadedNotes, err := notes.Load(conf.Notebooks)
		cobra.CheckErr(err)

		if cmd.Flags().Changed("tag") {
			tags, err := cmd.Flags().GetString("tag")
			cobra.CheckErr(err)
			loadedNotes = notes.Filter(loadedNotes, strings.Split(tags, ","))
		}

		out := ""
		for _, noteByName := range loadedNotes {
			for _, note := range notes.Sort(noteByName) {
				// If -n was provided, constrain the search only to matching
				// notes:
				if noteName != "" && note.Name != args[1] {
					continue
				}

				// Assume that we want to perform a case-insensitive search for
				// <phrase>, unless --regex is provided, in which case we pass
				// the regex unaltered:
				pattern := "(?i)" + phrase
				if regexFlag {
					pattern = phrase
				}

				reg, err := regexp.Compile(pattern)
				cobra.CheckErr(err)

				note.Body = note.Search(reg)
				if note.Body == "" {
					continue
				}

				if conf.Color() {
					note.Colorize(conf)
				}

				// Display the note body:
				out += fmt.Sprintf("%s %s\n%s\n",
					note.Name,
					display.Faint(fmt.Sprintf("(%s)", note.Notebook), conf),
					display.Indent(note.Body),
				)
			}
		}

		out = strings.TrimSpace(out)

		// NOTE: resist the temptation to call `display.Write` multiple times in the
		// loop above. That will not play nicely with the paginator.
		display.Write(out, conf)
	},
}
