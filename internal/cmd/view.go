package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/cheat/internal/cheatpath"
	"github.com/yagoyudi/cheat/internal/config"
	"github.com/yagoyudi/cheat/internal/display"
	"github.com/yagoyudi/cheat/internal/sheets"
)

func init() {
	viewCmd.Flags().BoolP("all", "A", false, "display cheatsheets from all cheatpaths")
	viewCmd.Flags().StringP("tag", "t", "", "filter cheatsheets by tag")
}

var viewCmd = &cobra.Command{
	Use:   "view [cheatsheet]",
	Short: "Displays a cheatsheet for viewing",
	Args:  cobra.ExactArgs(1),
	Example: `  cheat view kubectl
  cheat view kubectl -t community`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cheatsheet := args[0]

		allFlag, err := cmd.Flags().GetBool("all")
		if err != nil {
			return fmt.Errorf("cmd: %v", err)
		}
		tags, err := cmd.Flags().GetString("tag")
		if err != nil {
			return fmt.Errorf("cmd: %v", err)
		}

		var conf config.Config
		if err := viper.Unmarshal(&conf); err != nil {
			return fmt.Errorf("cmd: %v", err)
		}

		var cheatpaths []cheatpath.Cheatpath
		if err := viper.UnmarshalKey("cheatpaths", &cheatpaths); err != nil {
			return fmt.Errorf("cmd: %v", err)
		}

		// load the cheatsheets
		cheatsheets, err := sheets.Load(cheatpaths)
		if err != nil {
			return fmt.Errorf("cmd: %v", err)
		}

		// filter cheatcheats by tag if --tag was provided
		if cmd.Flags().Changed("tag") {
			cheatsheets = sheets.Filter(
				cheatsheets,
				strings.Split(tags, ","),
			)
		}

		// if --all was passed, display cheatsheets from all cheatpaths
		if allFlag {
			// iterate over the cheatpaths
			out := ""
			for _, cheatpath := range cheatsheets {

				// if the cheatpath contains the specified cheatsheet, display
				// it
				if sheet, ok := cheatpath[cheatsheet]; ok {

					// identify the matching cheatsheet
					out += fmt.Sprintf("%s %s\n",
						sheet.Title,
						display.Faint(fmt.Sprintf("(%s)", sheet.CheatPath), conf),
					)

					// apply colorization if requested
					if conf.Color() {
						sheet.Colorize(conf)
					}

					// display the cheatsheet
					out += display.Indent(sheet.Text) + "\n"
				}
			}

			// display and exit
			display.Write(strings.TrimSuffix(out, "\n"), conf)
			os.Exit(0)
		}

		// otherwise, consolidate the cheatsheets found on all paths into a single
		// map of `title` => `sheet` (ie, allow more local cheatsheets to override
		// less local cheatsheets)
		consolidated := sheets.Consolidate(cheatsheets)

		// fail early if the requested cheatsheet does not exist
		sheet, ok := consolidated[cheatsheet]
		if !ok {
			return fmt.Errorf("cmd: No cheatsheet found for '%s'\n", cheatsheet)
		}

		// apply colorization if requested
		if conf.Color() {
			sheet.Colorize(conf)
		}

		// display the cheatsheet
		display.Write(sheet.Text, conf)
		return nil
	},
}
