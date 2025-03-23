package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/cheat/internal/installer"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Args:    cobra.ExactArgs(0),
	Short:   `Setup cheat`,
	Example: "  cheat init",
	RunE: func(cmd *cobra.Command, _ []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "cheat: cmd: %v", err)
			os.Exit(1)
		}
		configPath := filepath.Join(home, ".config", "cheat")
		configFile := filepath.Join(configPath, "config.yaml")
		err = viper.ReadInConfig()
		if err != nil {
			_, ok := err.(viper.ConfigFileNotFoundError)
			if !ok {
				return err
			}
			// prompt the user to create a config file
			yes, err := installer.Prompt(
				"A config file was not found. Would you like to create one now? [Y/n]",
				true,
			)
			if err != nil {
				return err
			}

			// exit early on a negative answer
			if !yes {
				os.Exit(0)
			}

			// run the installer
			if err := installer.Run(configTemplate, configFile); err != nil {
				return err
			}

			// notify the user and exit
			fmt.Printf("Created config file: %s\n", configFile)
			fmt.Println("Please read this file for advanced configuration information.")
			fmt.Println()
		}
		fmt.Println(`All good to go!`)
		return nil
	},
}
