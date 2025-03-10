package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yagoyudi/cheat/internal/config"
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

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(viewCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(directoriesCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(searchCmd)

	rootCmd.Flags().BoolP("init", "i", false, "write a default config file to stdout")
}

var rootCmd = &cobra.Command{
	Use:   "cheat",
	Short: "Cheat allows you to create and view interactive cheatsheets on the command-line.",
	Long:  "It was designed to help remind *nix system administrators of options for commands that they use frequently, but not frequently enough to remember.",
	RunE: func(cmd *cobra.Command, _ []string) error {
		initFlag, err := cmd.Flags().GetBool("init")
		if err != nil {
			return err
		}

		if initFlag {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}

			// read the envvars into a map of strings
			envvars := map[string]string{}
			for _, e := range os.Environ() {
				pair := strings.SplitN(e, "=", 2)
				envvars[pair[0]] = pair[1]
			}

			// load the config template
			configs := configTemplate

			// identify the os-specifc paths at which configs may be located
			confpaths, err := config.Paths(runtime.GOOS, home, envvars)
			if err != nil {
				return err
			}

			// determine the appropriate paths for config data and (optional) community
			// cheatsheets based on the user's platform
			confpath := confpaths[0]
			confdir := filepath.Dir(confpath)

			// create paths for community and personal cheatsheets
			community := filepath.Join(confdir, "cheatsheets", "community")
			personal := filepath.Join(confdir, "cheatsheets", "personal")

			// template the above paths into the default configs
			configs = strings.Replace(configs, "COMMUNITY_PATH", community, -1)
			configs = strings.Replace(configs, "PERSONAL_PATH", personal, -1)

			// output the templated configs
			fmt.Println(configs)

			return nil
		}

		fmt.Println(cmd.Short + "\n" + cmd.Long + "\n")
		cmd.Usage()

		return nil
	},
}

func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		return fmt.Errorf("cmd: %v", err)
	}
	return nil
}
