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
	"github.com/yagoyudi/cheat/internal/config"
	"github.com/yagoyudi/cheat/internal/display"
	"github.com/yagoyudi/cheat/internal/sheet"
	"github.com/yagoyudi/cheat/internal/sheets"
)

func init() {
	listCmd.Flags().StringP("tag", "t", "", "filter cheatsheets by tag")
}

var listCmd = &cobra.Command{
	Use:   "ls [cheatsheet]",
	Short: "Lists all available cheatsheets",
	RunE: func(cmd *cobra.Command, args []string) error {
		tag, err := cmd.Flags().GetString("tag")
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
			return fmt.Errorf("cmd: failed to list cheatsheets: %v", err)
		}

		// filter cheatsheets by tag if --tag was provided
		if cmd.Flags().Changed("tag") {
			cheatsheets = sheets.Filter(
				cheatsheets,
				strings.Split(tag, ","),
			)
		}

		// instead of "consolidating" all of the cheatsheets (ie, overwriting
		// global sheets with local sheets), here we simply want to create a
		// slice containing all sheets.
		flattened := []sheet.Sheet{}
		for _, pathsheets := range cheatsheets {
			for _, s := range pathsheets {
				flattened = append(flattened, s)
			}
		}

		// sort the "flattened" sheets alphabetically
		sort.Slice(flattened, func(i, j int) bool {
			return flattened[i].Title < flattened[j].Title
		})

		// filter if <cheatsheet> was specified
		// NB: our docopt specification is misleading here. When used in conjunction
		// with `-l`, `<cheatsheet>` is really a pattern against which to filter
		// sheet titles.
		if len(args) >= 1 {
			cheatsheet := args[0]
			// initialize a slice of filtered sheets
			filtered := []sheet.Sheet{}

			// initialize our filter pattern
			pattern := "(?i)" + cheatsheet

			// compile the regex
			reg, err := regexp.Compile(pattern)
			if err != nil {
				return fmt.Errorf("cmd: failed to compile regexp: %s, %v", pattern, err)
			}

			// iterate over each cheatsheet, and pass-through those which match the
			// filter pattern
			for _, s := range flattened {
				if reg.MatchString(s.Title) {
					filtered = append(filtered, s)
				}
			}

			flattened = filtered
		}

		// return exit code 2 if no cheatsheets are available
		if len(flattened) == 0 {
			os.Exit(2)
		}

		// initialize a tabwriter to produce cleanly columnized output
		var out bytes.Buffer
		w := tabwriter.NewWriter(&out, 0, 0, 1, ' ', 0)

		// write a header row
		fmt.Fprintln(w, "title:\tfile:\ttags:")

		// generate sorted, columnized output
		for _, sheet := range flattened {
			fmt.Fprintf(
				w,
				"%s\t%s\t%s\n",
				sheet.Title, sheet.Path, strings.Join(sheet.Tags, ","),
			)
		}

		// write columnized output to stdout
		w.Flush()
		display.Write(out.String(), conf)

		return nil
	},
}
