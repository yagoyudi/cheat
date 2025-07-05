package cmd

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/note/internal/config"
	"github.com/yagoyudi/note/internal/display"
	"github.com/yagoyudi/note/internal/note"
	"github.com/yagoyudi/note/internal/notebook"
	"github.com/yagoyudi/note/internal/notes"
)

func init() {
	listCmd.Flags().StringP("tag", "t", "", "filter notes by tag")
	listCmd.Flags().StringP("book", "b", "", "filter notes by notebook")
}

var listCmd = &cobra.Command{
	Use:   "ls [note]",
	Short: "Lists all available notes",
	Example: `  note ls
  note ls -p personal
  note ls -t networking`,
	Run: func(cmd *cobra.Command, args []string) {
		var conf config.Config
		cobra.CheckErr(viper.Unmarshal(&conf))

		if cmd.Flags().Changed("book") {
			path, err := cmd.Flags().GetString("book")
			cobra.CheckErr(err)
			conf.Notebooks, err = notebook.Filter(conf.Notebooks, path)
			cobra.CheckErr(err)
		}

		loadedNotes, err := notes.Load(conf.Notebooks)
		cobra.CheckErr(err)

		if cmd.Flags().Changed("tag") {
			tag, err := cmd.Flags().GetString("tag")
			cobra.CheckErr(err)
			loadedNotes = notes.Filter(loadedNotes, strings.Split(tag, ","))
		}

		// Create a slice containing all notes:
		flattenedNotes := []note.Note{}
		for _, pathnotes := range loadedNotes {
			for _, note := range pathnotes {
				flattenedNotes = append(flattenedNotes, note)
			}
		}

		// Sort alphabetically:
		sort.Slice(flattenedNotes, func(i, j int) bool {
			return flattenedNotes[i].Name < flattenedNotes[j].Name
		})

		// Filter if [note] was specified:
		if len(args) >= 1 {
			noteName := args[0]
			filteredNotes := []note.Note{}

			pattern := "(?i)" + noteName
			reg, err := regexp.Compile(pattern)
			cobra.CheckErr(err)

			for _, note := range flattenedNotes {
				if reg.MatchString(note.Name) {
					filteredNotes = append(filteredNotes, note)
				}
			}
			flattenedNotes = filteredNotes
		}

		if len(flattenedNotes) == 0 {
			os.Exit(2)
		}

		var out bytes.Buffer
		w := tabwriter.NewWriter(&out, 0, 0, 1, ' ', 0)
		fmt.Fprintln(w, "title:\tfile:\ttags:")
		for _, note := range flattenedNotes {
			fmt.Fprintf(w, "%s\t%s\t%s\n", note.Name, note.Path, strings.Join(note.Tags, ","))
		}
		w.Flush()
		display.Write(out.String(), conf)
	},
}
