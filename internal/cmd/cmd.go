package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/cheat/internal/installer"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cheat: cmd: %v", err)
		os.Exit(1)
	}
	configPath := filepath.Join(home, ".config", "cheat")
	configFile := filepath.Join(configPath, "config.yaml")
	viper.AddConfigPath(configPath)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// prompt the user to create a config file
			yes, err := installer.Prompt(
				"A config file was not found. Would you like to create one now? [Y/n]",
				true,
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "cheat: cmd: %v\n", err)
				os.Exit(1)
			}

			// exit early on a negative answer
			if !yes {
				os.Exit(0)
			}

			// run the installer
			if err := installer.Run(configTemplate, configFile); err != nil {
				fmt.Fprintf(os.Stderr, "cheat: cmd: %v\n", err)
				os.Exit(1)
			}

			// notify the user and exit
			fmt.Printf("Created config file: %s\n", configFile)
			fmt.Println("Please read this file for advanced configuration information.")
			fmt.Println()
		} else {
			fmt.Fprintf(os.Stderr, "cheat: cmd: %v\n", err)
			os.Exit(1)
		}
	}

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
