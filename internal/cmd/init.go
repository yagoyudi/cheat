package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/note/internal/installer"
)

const configTemplate = `---
# The editor to use with 'note e <sheet>'. Defaults to $EDITOR or $VISUAL.
editor: EDITOR_PATH

# Should 'note' always colorize output?
colorize: true

# Which 'chroma' colorscheme should be applied to the output?
# Options are available here:
#   https://github.com/alecthomas/chroma/tree/master/styles
style: monokai

# Which 'chroma' "formatter" should be applied?
# One of: "terminal", "terminal256", "terminal16m"
formatter: terminal256

# Through which pager should output be piped?
# 'less -FRX' is recommended on Unix systems
# 'more' is recommended on Windows
pager: PAGER_PATH

# Notebooks are paths at which notes are available on your local # filesystem.
#
# It is useful to sort notes into different notes for organizational
# purposes. For example, you might want one note for community
# notes, one for personal notes, one for notes pertaining to
# your day job, one for code snippets, etc.
#
# notes are scoped, such that more "local" notes take priority over
# more "global" notes. (The most global note is listed first in this
# file; the most local is listed last.) For example, if there is a 'tar'
# note on both global and local paths, you'll be presented with the local
# one by default. ('note -p' can be used to view notes from alternative
# notes.)
#
# notes can also be tagged as "read only". This instructs note not to
# automatically create notes on a read-only note. Instead, when you
# would like to edit a read-only note using 'note e', note will
# perform a copy-on-write of that note from a read-only note to a
# writeable note.
#
# This is very useful when you would like to maintain, for example, a
# "pristine" repository of community notes on one note, and an
# editable personal reponsity of notes on another note.
#
# notes can be also configured to automatically apply tags to notes
# on certain paths, which can be useful for querying purposes.
# Example: 'note -t work jenkins'.
#
# Community notes must be installed separately, though you may have
# downloaded them automatically when installing 'note'. If not, you may
# download them here:
#
# https://github.com/yagoyudi/note
notebook:
  # note properties mean the following:
  #   'name': the name of the note (view with 'note -d', filter with 'note -p')
  #   'path': the filesystem path of the note directory (view with 'note -d')
  #   'tags': tags that should be automatically applied to sheets on this path
  #   'readonly': shall user-created ('note -e') notes be saved here?
  - name: community
    path: COMMUNITY_PATH
    tags: [ community ]
    readonly: true

  # If you have personalized notes, list them last. They will take
  # precedence over the more global notes.
  - name: personal
    path: PERSONAL_PATH
    tags: [ personal ]
    readonly: false

  # While it requires no configuration here, it's also worth noting that
  # note will automatically append directories named '.note' within the
  # current working directory to the 'note'. This can be very useful if
  # you'd like to closely associate notes with, for example, a directory
  # containing source code.
  #
  # Such "directory-scoped" notes will be treated as the most "local"
  # notes, and will override less "local" notes. Similarly,
  # directory-scoped notes will always be editable ('readonly: false').
`

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Args:    cobra.ExactArgs(0),
	Short:   `Setup note`,
	Example: "  note init",
	Run: func(cmd *cobra.Command, _ []string) {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		configPath := filepath.Join(home, ".config", "note")
		configFile := filepath.Join(configPath, "config.yaml")
		err = viper.ReadInConfig()
		if err != nil {
			_, ok := err.(viper.ConfigFileNotFoundError)
			if !ok {
				cobra.CheckErr(err)
			}
			yes, err := installer.Prompt("A config file was not found. Would you like to create one now? [Y/n]", true)
			cobra.CheckErr(err)

			if !yes {
				os.Exit(0)
			}
			cobra.CheckErr(installer.Run(configTemplate, configFile))

			fmt.Printf("Created config file: %s\n", configFile)
			fmt.Println("Please read this file for advanced configuration information.")
			fmt.Println()
		}
		fmt.Println("All good to go!")
	},
}
