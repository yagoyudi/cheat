package cmd

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/note/internal/config"
	"github.com/yagoyudi/note/internal/display"
)

var notebooksCmd = &cobra.Command{
	Use:     "books",
	Aliases: []string{"b"},
	Short:   "Lists the configured notebooks",
	Example: "  note books",
	Run: func(cmd *cobra.Command, _ []string) {
		var conf config.Config
		cobra.CheckErr(viper.Unmarshal(&conf))

		var out bytes.Buffer
		w := tabwriter.NewWriter(&out, 0, 0, 1, ' ', 0)

		for _, notebook := range conf.Notebooks {
			fmt.Fprintf(w, "%s:\t%s\n", notebook.Name, notebook.Path)
		}

		w.Flush()
		display.Write(out.String(), conf)
	},
}
