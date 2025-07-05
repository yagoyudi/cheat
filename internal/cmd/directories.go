package cmd

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/cheat/internal/config"
	"github.com/yagoyudi/cheat/internal/display"
)

var notebooksCmd = &cobra.Command{
	Use:     "notebooks",
	Aliases: []string{"nb"},
	Short:   "Lists the configured notebooks",
	Example: "  cheat nb",
	Run: func(cmd *cobra.Command, _ []string) {
		var conf config.Config
		cobra.CheckErr(viper.Unmarshal(&conf))

		// Initialize a tabwriter to produce cleanly columnized output:
		var out bytes.Buffer
		w := tabwriter.NewWriter(&out, 0, 0, 1, ' ', 0)

		// Generate sorted, columnized output:
		for _, notebook := range conf.Notebooks {
			fmt.Fprintf(w, "%s:\t%s\n", notebook.Name, notebook.Path)
		}

		w.Flush()
		display.Write(out.String(), conf)
	},
}
