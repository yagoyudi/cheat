// Package config implements functions pertaining to configuration management.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yagoyudi/cheat/internal/cheatpath"

	"gopkg.in/yaml.v3"
)

// Config encapsulates configuration parameters
type Config struct {
	Colorize   bool                  `yaml:"colorize"`
	Editor     string                `yaml:"editor"`
	Cheatpaths []cheatpath.Cheatpath `yaml:"cheatpaths"`
	Style      string                `yaml:"style"`
	Formatter  string                `yaml:"formatter"`
	Pager      string                `yaml:"pager"`
	Path       string
}

// New returns a new Config struct
func New(_ map[string]interface{}, confPath string, resolve bool) (Config, error) {

	// read the config file
	buf, err := os.ReadFile(confPath)
	if err != nil {
		return Config{}, fmt.Errorf("config: could not read config file: %v", err)
	}

	// initialize a config object
	conf := Config{}

	// store the config path
	conf.Path = confPath

	// unmarshal the yaml
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return Config{}, fmt.Errorf("config: could not unmarshal yaml: %v", err)
	}

	// if a .cheat directory exists locally, append it to the cheatpaths
	cwd, err := os.Getwd()
	if err != nil {
		return Config{}, fmt.Errorf("config: failed to get cwd: %v", err)
	}

	local := filepath.Join(cwd, ".cheat")
	if _, err := os.Stat(local); err == nil {
		path := cheatpath.Cheatpath{
			Name:     "cwd",
			Path:     local,
			ReadOnly: false,
			Tags:     []string{},
		}

		conf.Cheatpaths = append(conf.Cheatpaths, path)
	}

	// process cheatpaths
	for i, cheatpath := range conf.Cheatpaths {

		// expand ~ in config paths
		expanded, err := expandPath(cheatpath.Path)
		if err != nil {
			return Config{}, fmt.Errorf("config: failed to expand ~: %v", err)
		}

		// follow symlinks
		//
		// NB: `resolve` is an ugly kludge that exists for the sake of unit-tests.
		// It's necessary because `EvalSymlinks` will error if the symlink points
		// to a non-existent location on the filesystem. When unit-testing,
		// however, we don't want to have dependencies on the filesystem. As such,
		// `resolve` is a switch that allows us to turn off symlink resolution when
		// running the config tests.
		if resolve {
			evaled, err := filepath.EvalSymlinks(expanded)
			if err != nil {
				return Config{}, fmt.Errorf(
					"config: failed to resolve symlink: %s: %v",
					expanded,
					err,
				)
			}

			expanded = evaled
		}

		conf.Cheatpaths[i].Path = expanded
	}

	// if an editor was not provided in the configs, attempt to choose one
	// that's appropriate for the environment
	if conf.Editor == "" {
		conf.Editor, err = Editor()
		if err != nil {
			return Config{}, err
		}
	}

	// if a chroma style was not provided, set a default
	if conf.Style == "" {
		conf.Style = "bw"
	}

	// if a chroma formatter was not provided, set a default
	if conf.Formatter == "" {
		conf.Formatter = "terminal"
	}

	// load the pager
	conf.Pager = strings.TrimSpace(conf.Pager)

	return conf, nil
}
