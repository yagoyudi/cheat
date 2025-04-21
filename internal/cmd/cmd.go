package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cheat: cmd: %v", err)
		os.Exit(1)
	}
	configPath := filepath.Join(home, ".config", "cheat", "config.yaml")
	viper.AddConfigPath(configPath)
	viper.AutomaticEnv()

	rootCmd.AddCommand(
		listCmd,
		viewCmd,
		editCmd,
		versionCmd,
		directoriesCmd,
		removeCmd,
		searchCmd,
		initCmd,
		tagsCmd,
	)
}

var rootCmd = &cobra.Command{
	Use:   "cheat",
	Short: "Cheat allows you to create and view interactive cheatsheets on the command-line.",
	Long:  "Cheat was designed to help remind *nix system administrators of options for commands that they use frequently, but not frequently enough to remember.",
}

func Execute() {
	rootCmd.Execute()
}
