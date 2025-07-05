package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/cheat/internal/config"
	"github.com/yagoyudi/cheat/internal/notes"
)

func init() {
	removeCmd.Flags().StringP("tag", "t", "", "filter cheatsheets by tag")
}

var removeCmd = &cobra.Command{
	Use:     "rm [cheatsheet]",
	Short:   "Removes a cheatsheet",
	Args:    cobra.ExactArgs(1),
	Example: `  cheat rm kubectl -t community`,
	RunE: func(cmd *cobra.Command, args []string) error {
		tags, err := cmd.Flags().GetString("tag")
		if err != nil {
			return err
		}

		cheatsheet := args[0]

		var conf config.Config
		if err := viper.Unmarshal(&conf); err != nil {
			return err
		}

		// load the cheatsheets
		cheatsheets, err := notes.Load(conf.Notebooks)
		if err != nil {
			return err
		}

		// filter cheatcheats by tag if --tag was provided
		if cmd.Flags().Changed("tag") {
			cheatsheets = notes.Filter(
				cheatsheets,
				strings.Split(tags, ","),
			)
		}

		// consolidate the cheatsheets found on all paths into a single map of
		// `title` => `sheet` (ie, allow more local cheatsheets to override less
		// local cheatsheets)
		consolidated := notes.Consolidate(cheatsheets)

		// fail early if the requested cheatsheet does not exist
		sheet, ok := consolidated[cheatsheet]
		if !ok {
			return fmt.Errorf("No cheatsheet found for '%s'.\n", cheatsheet)
		}

		// fail early if the sheet is read-only
		if sheet.ReadOnly {
			return fmt.Errorf("cheatsheet '%s' is read-only.\n", cheatsheet)
		}

		// otherwise, attempt to delete the sheet
		if err := os.Remove(sheet.Path); err != nil {
			return err
		}

		return nil
	},
}
