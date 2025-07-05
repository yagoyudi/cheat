// Implements functions pertaining to configuration management
package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/yagoyudi/note/internal/notebook"

	"gopkg.in/yaml.v3"
)

// Encapsulates configuration parameters
type Config struct {
	Colorize  bool                `yaml:"colorize"`
	Editor    string              `yaml:"editor"`
	Notebooks []notebook.Notebook `yaml:"notepaths"`
	Style     string              `yaml:"style"`
	Formatter string              `yaml:"formatter"`
	Pager     string              `yaml:"pager"`
	Path      string
}

// Returns a new Config struct
func New(_ map[string]any, confPath string, resolve bool) (Config, error) {
	buf, err := os.ReadFile(confPath)
	if err != nil {
		return Config{}, fmt.Errorf("config: could not read config file: %v", err)
	}

	conf := Config{}
	conf.Path = confPath

	if err := yaml.Unmarshal(buf, &conf); err != nil {
		return Config{}, fmt.Errorf("config: could not unmarshal yaml: %v", err)
	}

	// If a .cheat directory exists locally, append it to the notepaths:
	cwd, err := os.Getwd()
	if err != nil {
		return Config{}, fmt.Errorf("config: failed to get cwd: %v", err)
	}

	local := filepath.Join(cwd, ".cheat")
	if _, err := os.Stat(local); err == nil {
		path := notebook.Notebook{
			Name:     "cwd",
			Path:     local,
			ReadOnly: false,
			Tags:     []string{},
		}

		conf.Notebooks = append(conf.Notebooks, path)
	}

	for i, notebook := range conf.Notebooks {
		expanded, err := expandPath(notebook.Path)
		if err != nil {
			return Config{}, fmt.Errorf("config: failed to expand ~: %v", err)
		}

		// Follow symlinks:
		//
		// NOTE: `resolve` is an ugly kludge that exists for the sake of
		// unit-tests. It's necessary because `EvalSymlinks` will error if the
		// symlink points to a non-existent location on the filesystem. When
		// unit-testing, however, we don't want to have dependencies on the
		// filesystem. As such, `resolve` is a switch that allows us to turn
		// off symlink resolution when running the config tests.
		if resolve {
			evaled, err := filepath.EvalSymlinks(expanded)
			if err != nil {
				return Config{}, fmt.Errorf("config: failed to resolve symlink: %s: %v", expanded, err)
			}
			expanded = evaled
		}
		conf.Notebooks[i].Path = expanded
	}

	if conf.Editor == "" {
		conf.Editor, err = Editor()
		if err != nil {
			return Config{}, err
		}
	}
	if conf.Style == "" {
		conf.Style = "bw"
	}
	if conf.Formatter == "" {
		conf.Formatter = "terminal"
	}
	conf.Pager = strings.TrimSpace(conf.Pager)

	return conf, nil
}

// Attempts to locate an editor that's appropriate for the environment
func Editor() (string, error) {
	def, _ := exec.LookPath("editor")
	nano, _ := exec.LookPath("nano")
	vim, _ := exec.LookPath("vim")
	editors := []string{
		os.Getenv("VISUAL"),
		os.Getenv("EDITOR"),
		def,
		nano,
		vim,
	}

	// Return the first editor that was found per the priority above:
	for _, editor := range editors {
		if editor != "" {
			return editor, nil
		}
	}
	return "", fmt.Errorf("config: no editor set")
}

// Indicates whether colorization should be applied to the output
func (c *Config) Color() bool {
	colorize := c.Colorize
	// Only apply colorization if we're writing to a tty:
	if !isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		colorize = false
	}
	return colorize
}

// Initializes a config file
func Init(confpath string, configs string) error {
	// Assert that the config directory exists:
	if err := os.MkdirAll(filepath.Dir(confpath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write the config file:
	if err := os.WriteFile(confpath, []byte(configs), 0644); err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	return nil
}

// Attempts to locate a pager that's appropriate for the environment.
func Pager() string {
	if os.Getenv("PAGER") != "" {
		return os.Getenv("PAGER")
	}

	// Search for `pager`, `less`, and `more` on the `$PATH`. If none are
	// found, return an empty pager.
	for _, pager := range []string{"pager", "less", "more"} {
		if path, err := exec.LookPath(pager); err != nil {
			return path
		}
	}
	return ""
}

// Returns the config file path
func Path(paths []string) (string, error) {
	// Check if the config file exists on any paths:
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}
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

// Returns config file paths that are appropriate for the operating system
func Paths(sys string, home string, envvars map[string]string) ([]string, error) {
	// if `CHEAT_CONFIG_PATH` is set, expand ~ and return it
	if confpath, ok := envvars["CHEAT_CONFIG_PATH"]; ok {
		expanded, err := expandPath(confpath)
		if err != nil {
			return []string{}, fmt.Errorf("failed to expand ~: %v", err)
		}

		return []string{expanded}, nil
	}

	switch sys {
	case "aix", "android", "darwin", "dragonfly", "freebsd", "illumos", "ios",
		"linux", "netbsd", "openbsd", "plan9", "solaris":
		paths := []string{}

		// Don't include the `XDG_CONFIG_HOME` path if that envvar is not set:
		if xdgpath, ok := envvars["XDG_CONFIG_HOME"]; ok {
			paths = append(paths, filepath.Join(xdgpath, "note", "conf.yml"))
		}

		paths = append(paths, []string{
			filepath.Join(home, ".config", "note", "conf.yml"),
			filepath.Join(home, ".note", "conf.yml"),
			"/etc/note/conf.yml",
		}...)

		return paths, nil
	case "windows":
		return []string{
			filepath.Join(envvars["APPDATA"], "note", "conf.yml"),
			filepath.Join(envvars["PROGRAMDATA"], "note", "conf.yml"),
		}, nil
	default:
		return []string{}, fmt.Errorf("unsupported os: %s", sys)
	}
}

// Returns an error if the config is invalid
func (c *Config) Validate() error {
	if c.Editor == "" {
		return fmt.Errorf("config error: editor unspecified")
	}

	// Assert that at least one cheatpath was specified:
	if len(c.Notebooks) == 0 {
		return fmt.Errorf("config error: no notebooks specified")
	}

	// Assert that each path and name is unique:
	names := make(map[string]bool)
	paths := make(map[string]bool)
	for _, cheatpath := range c.Notebooks {
		if err := cheatpath.Validate(); err != nil {
			return fmt.Errorf("config error: %v", err)
		}

		if _, ok := names[cheatpath.Name]; ok {
			return fmt.Errorf(
				"config error: notebook name is not unique: %s",
				cheatpath.Name,
			)
		}
		names[cheatpath.Name] = true

		if _, ok := paths[cheatpath.Path]; ok {
			return fmt.Errorf(
				"config error: notebook path is not unique: %s",
				cheatpath.Path,
			)
		}
		paths[cheatpath.Path] = true
	}

	// TODO: assert valid styles?

	// Assert that the formatter is valid:
	formatters := map[string]bool{
		"terminal":    true,
		"terminal256": true,
		"terminal16m": true,
	}
	if _, ok := formatters[c.Formatter]; !ok {
		return fmt.Errorf("config error: header is invalid: %s", c.Formatter)
	}

	return nil
}
