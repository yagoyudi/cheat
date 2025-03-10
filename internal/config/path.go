package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Path returns the config file path
func Path(paths []string) (string, error) {

	// check if the config file exists on any paths
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	// we can't find the config file if we make it this far
	return "", fmt.Errorf("could not locate config file")
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %v", err)
		}
		return filepath.Join(home, path[1:]), nil
	}
	return path, nil
}
