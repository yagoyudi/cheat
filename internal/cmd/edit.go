package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/cheat/internal/notebook"
	"github.com/yagoyudi/cheat/internal/notes"
)

func init() {
	editCmd.Flags().StringP("tag", "t", "", "filter cheatsheets by tag")
}

var editCmd = &cobra.Command{
	Use:     "edit [cheatsheet]",
	Aliases: []string{"e"},
	Short:   "Opens a cheatsheet for editing (or creates it if it doesn't exist)",
	Args:    cobra.ExactArgs(1),
	Example: `  cheat edit tar     # opens the "tar" cheatsheet for editing, or creates it if it does not exist
  cheat edit foo/bar # nested cheatsheets are accessed like this`,
	Run: func(cmd *cobra.Command, args []string) {
		var notebooks []notebook.Notebook
		cobra.CheckErr(viper.UnmarshalKey("cheatpaths", &notebooks))

		loadedNotes, err := notes.Load(notebooks)
		cobra.CheckErr(err)

		tag, err := cmd.Flags().GetString("tag")
		cobra.CheckErr(err)

		// Filter notes by tag if --tag was provided:
		if cmd.Flags().Changed("tag") {
			loadedNotes = notes.Filter(loadedNotes, strings.Split(tag, ","))
		}

		// consolidate the cheatsheets found on all paths into a single map of
		// `title` => `sheet` (ie, allow more local cheatsheets to override less
		// local cheatsheets)
		consolidated := notes.Consolidate(loadedNotes)

		// The file path of the note to edit:
		var editpath string

		// Determine if the sheet exists:
		noteName := args[0]
		note, ok := consolidated[noteName]
		if ok {
			if note.ReadOnly {
				// compute the new edit path
				// begin by getting a writeable cheatpath
				writepath, err := notebook.Writeable(notebooks)
				cobra.CheckErr(err)

				// compute the new edit path
				editpath = filepath.Join(writepath.Path, note.Name)

				// create any necessary subdirectories
				dirs := filepath.Dir(editpath)
				if dirs != "." {
					cobra.CheckErr(os.MkdirAll(dirs, 0755))
				}

				// Copy the sheet to the new edit path:
				cobra.CheckErr(note.Copy(editpath))

			} else {
				editpath = note.Path
			}
		} else {
			// Create note:

			// compute the new edit path
			// begin by getting a writeable cheatpath
			writepath, err := notebook.Writeable(notebooks)
			cobra.CheckErr(err)

			// compute the new edit path
			editpath = filepath.Join(writepath.Path, noteName)

			// create any necessary subdirectories
			dirs := filepath.Dir(editpath)
			if dirs != "." {
				cobra.CheckErr(os.MkdirAll(dirs, 0755))
			}
		}

		parts := strings.Fields(viper.GetString("editor"))
		editor := parts[0]
		editorArgs := append(parts[1:], editpath)
		c := exec.Command(editor, editorArgs...)
		c.Stdout = os.Stdout
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		cobra.CheckErr(c.Run())
	},
}
