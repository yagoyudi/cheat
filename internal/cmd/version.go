package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "0.1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display binary version in use",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("v" + version)
	},
}
