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

var directoriesCmd = &cobra.Command{
	Use:     "dirs",
	Short:   "Lists the configured cheatpaths",
	Example: "  cheat dirs",
	RunE: func(cmd *cobra.Command, _ []string) error {
		var conf config.Config
		if err := viper.Unmarshal(&conf); err != nil {
			return err
		}

		// initialize a tabwriter to produce cleanly columnized output
		var out bytes.Buffer
		w := tabwriter.NewWriter(&out, 0, 0, 1, ' ', 0)

		// generate sorted, columnized output
		for _, path := range conf.Cheatpaths {
			fmt.Fprintf(w, "%s:\t%s\n", path.Name, path.Path)
		}

		// write columnized output to stdout
		w.Flush()
		display.Write(out.String(), conf)
		return nil
	},
}
