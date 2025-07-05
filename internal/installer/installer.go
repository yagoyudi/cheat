// Implements functions that provide a first-time installation wizard.
package installer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yagoyudi/cheat/internal/config"
	"github.com/yagoyudi/cheat/internal/repo"
)

// Prompts the user for a answer
func Prompt(prompt string, def bool) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s: ", prompt)

	ans, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("installer: %v", err)
	}

	ans = strings.ToLower(strings.TrimSpace(ans))
	switch ans {
	case "y":
		return true, nil
	case "":
		return def, nil
	default:
		return false, nil
	}
}

// Runs the installer
func Run(configs string, confpath string) error {
	// Determine the appropriate paths for config data and (optional) community
	// cheatsheets based on the user's platform:
	confdir := filepath.Dir(confpath)

	// Create paths for community and personal cheatsheets:
	community := filepath.Join(confdir, "cheatsheets", "community")
	personal := filepath.Join(confdir, "cheatsheets", "personal")

	// Set default cheatpaths:
	configs = strings.Replace(configs, "COMMUNITY_PATH", community, -1)
	configs = strings.Replace(configs, "PERSONAL_PATH", personal, -1)

	// Locate and set a default pager:
	configs = strings.Replace(configs, "PAGER_PATH", config.Pager(), -1)

	// Locate and set a default editor:
	if editor, err := config.Editor(); err == nil {
		configs = strings.Replace(configs, "EDITOR_PATH", editor, -1)
	}

	// Prompt the user to download the community cheatsheets:
	yes, err := Prompt("Would you like to download the community cheatsheets? [Y/n]", true)
	if err != nil {
		return fmt.Errorf("failed to prompt: %v", err)
	}

	// Clone the community cheatsheets if so instructed:
	if yes {
		fmt.Printf("Cloning community cheatsheets to %s.\n", community)
		if err := repo.Clone(community); err != nil {
			return fmt.Errorf("failed to clone cheatsheets: %v", err)
		}

		// Create a directory for personal cheatsheets:
		fmt.Printf("Cloning personal cheatsheets to %s.\n", personal)
		if err := os.MkdirAll(personal, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	// If the config file does not exist, try to create one:
	if err = config.Init(confpath, configs); err != nil {
		return fmt.Errorf("failed to create config file: %v", err)
	}

	return nil
}
