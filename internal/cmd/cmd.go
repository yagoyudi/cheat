package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	configPath := filepath.Join(home, ".config", "note")
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	cobra.CheckErr(viper.ReadInConfig())

	rootCmd.AddCommand(
		listCmd,
		viewCmd,
		editCmd,
		notebooksCmd,
		removeCmd,
		searchCmd,
		initCmd,
		tagsCmd,
	)
}

var rootCmd = &cobra.Command{
	Use:   "note",
	Short: "Note allows you to create and view interactive notes on the command-line.",
	Long:  "Note was designed to help remind *nix system administrators of options for commands that they use frequently, but not frequently enough to remember.",
}

func Execute() {
	rootCmd.Execute()
}
