package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/cheat/internal/cheatpath"
	"github.com/yagoyudi/cheat/internal/sheets"
)

func init() {
	editCmd.Flags().StringP("tag", "t", "", "filter cheatsheets by tag")
}

var editCmd = &cobra.Command{
	Use:   "edit [cheatsheet]",
	Short: "Opens a cheatsheet for editing (or creates it if it doesn't exist)",
	Args:  cobra.ExactArgs(1),
	Example: `  cheat edit tar     # opens the "tar" cheatsheet for editing, or creates it if it does not exist
  cheat edit foo/bar # nested cheatsheets are accessed like this`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var cheatpaths []cheatpath.Cheatpath
		if err := viper.UnmarshalKey("cheatpaths", &cheatpaths); err != nil {
			return err
		}

		// load the cheatsheets
		cheatsheets, err := sheets.Load(cheatpaths)
		if err != nil {
			return err
		}

		tag, err := cmd.Flags().GetString("tag")
		if err != nil {
			return err
		}

		// filter cheatsheets by tag if --tag was provided
		if cmd.Flags().Changed("tag") {
			cheatsheets = sheets.Filter(
				cheatsheets,
				strings.Split(tag, ","),
			)
		}

		// consolidate the cheatsheets found on all paths into a single map of
		// `title` => `sheet` (ie, allow more local cheatsheets to override less
		// local cheatsheets)
		consolidated := sheets.Consolidate(cheatsheets)

		// the file path of the sheet to edit
		var editpath string

		// determine if the sheet exists
		cheatsheet := args[0]
		sheet, ok := consolidated[cheatsheet]

		// if the sheet exists and is not read-only, edit it in place
		if ok && !sheet.ReadOnly {
			editpath = sheet.Path

			// if the sheet exists but is read-only, copy it before editing
		} else if ok && sheet.ReadOnly {
			// compute the new edit path
			// begin by getting a writeable cheatpath
			writepath, err := cheatpath.Writeable(cheatpaths)
			if err != nil {
				return err
			}

			// compute the new edit path
			editpath = filepath.Join(writepath.Path, sheet.Title)

			// create any necessary subdirectories
			dirs := filepath.Dir(editpath)
			if dirs != "." {
				if err := os.MkdirAll(dirs, 0755); err != nil {
					return err
				}
			}

			// copy the sheet to the new edit path
			err = sheet.Copy(editpath)
			if err != nil {
				return err
			}

			// if the sheet does not exist, create it
		} else {
			// compute the new edit path
			// begin by getting a writeable cheatpath
			writepath, err := cheatpath.Writeable(cheatpaths)
			if err != nil {
				return err
			}

			// compute the new edit path
			editpath = filepath.Join(writepath.Path, cheatsheet)

			// create any necessary subdirectories
			dirs := filepath.Dir(editpath)
			if dirs != "." {
				if err := os.MkdirAll(dirs, 0755); err != nil {
					return err
				}
			}
		}

		// split `conf.Editor` into parts to separate the editor's executable
		// from any arguments it may have been passed. If this is not done, the
		// nearby call to `exec.Command` will fail.
		parts := strings.Fields(viper.GetString("editor"))
		editor := parts[0]
		editorArgs := append(parts[1:], editpath)

		c := exec.Command(editor, editorArgs...)
		c.Stdout = os.Stdout
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		return c.Run()
	},
}
