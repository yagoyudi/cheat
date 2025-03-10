package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/cheat/internal/config"
	"github.com/yagoyudi/cheat/internal/display"
	"github.com/yagoyudi/cheat/internal/sheets"
)

func init() {
	searchCmd.Flags().StringP("tag", "t", "", "filter cheatsheets by tag")
	searchCmd.Flags().StringP("cheatsheet", "c", "", "constrain the search only to matching cheatsheets")
	searchCmd.Flags().BoolP("regex", "r", false, "treat search <phrase> as a regex")
}

var searchCmd = &cobra.Command{
	Use:   "search [phrase]",
	Short: "Searches for strings in cheatsheets",
	RunE: func(cmd *cobra.Command, args []string) error {
		phrase := args[0]

		tags, err := cmd.Flags().GetString("tag")
		if err != nil {
			return err
		}
		regexFlag, err := cmd.Flags().GetBool("regex")
		if err != nil {
			return err
		}
		cheatsheet, err := cmd.Flags().GetString("cheatsheet")
		if err != nil {
			return err
		}

		var conf config.Config
		if err := viper.Unmarshal(&conf); err != nil {
			return err
		}

		// load the cheatsheets
		cheatsheets, err := sheets.Load(conf.Cheatpaths)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to list cheatsheets: %v\n", err)
			os.Exit(1)
		}

		// filter cheatcheats by tag if --tag was provided
		if cmd.Flags().Changed("tag") {
			cheatsheets = sheets.Filter(
				cheatsheets,
				strings.Split(tags, ","),
			)
		}

		// iterate over each cheatpath
		out := ""
		for _, pathcheats := range cheatsheets {

			// sort the cheatsheets alphabetically, and search for matches
			for _, sheet := range sheets.Sort(pathcheats) {

				// if -c was provided, constrain the search only to matching
				// cheatsheets
				if cheatsheet != "" && sheet.Title != args[1] {
					continue
				}

				// assume that we want to perform a case-insensitive search for <phrase>
				pattern := "(?i)" + phrase

				// unless --regex is provided, in which case we pass the regex unaltered
				if regexFlag {
					pattern = phrase
				}

				// compile the regex
				reg, err := regexp.Compile(pattern)
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to compile regexp: %s, %v\n", pattern, err)
					os.Exit(1)
				}

				// `Search` will return text entries that match the search terms.
				// We're using it here to overwrite the prior cheatsheet Text,
				// filtering it to only what is relevant.
				sheet.Text = sheet.Search(reg)

				// if the sheet did not match the search, ignore it and move on
				if sheet.Text == "" {
					continue
				}

				// if colorization was requested, apply it here
				if conf.Color() {
					sheet.Colorize(conf)
				}

				// display the cheatsheet body
				out += fmt.Sprintf(
					"%s %s\n%s\n",
					// append the cheatsheet title
					sheet.Title,
					// append the cheatsheet path
					display.Faint(fmt.Sprintf("(%s)", sheet.CheatPath), conf),
					// indent each line of content
					display.Indent(sheet.Text),
				)
			}
		}

		// trim superfluous newlines
		out = strings.TrimSpace(out)

		// display the output
		// NB: resist the temptation to call `display.Write` multiple times in the
		// loop above. That will not play nicely with the paginator.
		display.Write(out, conf)
		return nil
	},
}
