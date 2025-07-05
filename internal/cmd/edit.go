package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/note/internal/notebook"
	"github.com/yagoyudi/note/internal/notes"
)

func init() {
	editCmd.Flags().StringP("tag", "t", "", "filter notes by tag")
}

var editCmd = &cobra.Command{
	Use:     "edit [note]",
	Aliases: []string{"e"},
	Short:   "Opens a note for editing (or creates it if it doesn't exist)",
	Args:    cobra.ExactArgs(1),
	Example: `  note edit tar
  note e tar
	`,
	Run: func(cmd *cobra.Command, args []string) {
		var notebooks []notebook.Notebook
		cobra.CheckErr(viper.UnmarshalKey("notebooks", &notebooks))

		loadedNotes, err := notes.Load(notebooks)
		cobra.CheckErr(err)

		tag, err := cmd.Flags().GetString("tag")
		cobra.CheckErr(err)

		// Filter notes by tag if --tag was provided:
		if cmd.Flags().Changed("tag") {
			loadedNotes = notes.Filter(loadedNotes, strings.Split(tag, ","))
		}

		consolidatedNotes := notes.Consolidate(loadedNotes)

		// The file path of the note to edit:
		var editpath string

		// Determine if the note exists:
		noteName := args[0]
		note, ok := consolidatedNotes[noteName]
		if ok {
			if note.ReadOnly {
				// Compute the new edit path:
				writepath, err := notebook.Writeable(notebooks)
				cobra.CheckErr(err)
				editpath = filepath.Join(writepath.Path, note.Name)
				dirs := filepath.Dir(editpath)
				if dirs != "." {
					cobra.CheckErr(os.MkdirAll(dirs, 0755))
				}

				// Copy the note to the new edit path:
				cobra.CheckErr(note.Copy(editpath))
			} else {
				editpath = note.Path
			}
		} else {
			// Create note:
			writepath, err := notebook.Writeable(notebooks)
			cobra.CheckErr(err)
			editpath = filepath.Join(writepath.Path, noteName)
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
