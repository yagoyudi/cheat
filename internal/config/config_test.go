package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yagoyudi/cheat/internal/mock"
	"github.com/yagoyudi/cheat/internal/notebook"
)

// TestConfig asserts that the configs are loaded correctly
func TestConfigSuccessful(t *testing.T) {

	// initialize a config
	conf, err := New(map[string]interface{}{}, mock.Path("conf/conf.yml"), false)
	assert.NoError(t, err, "failed to parse config file")
	assert.Equal(t, "vim", conf.Editor, "failed to set editor")
	assert.Equal(t, true, conf.Colorize, "failed to set colorize")

	// get the user's home directory (with ~ expanded)
	home, err := os.UserHomeDir()
	assert.NoError(t, err, "failed to get homedir")

	// assert that the cheatpaths are correct
	want := []notebook.Notebook{
		{
			Path:     filepath.Join(home, ".dotfiles", "cheat", "community"),
			ReadOnly: true,
			Tags:     []string{"community"},
		},
		{
			Path:     filepath.Join(home, ".dotfiles", "cheat", "work"),
			ReadOnly: false,
			Tags:     []string{"work"},
		},
		{
			Path:     filepath.Join(home, ".dotfiles", "cheat", "personal"),
			ReadOnly: false,
			Tags:     []string{"personal"},
		},
	}
	assert.Equal(t, want, conf.Notebooks, "failed to return expected results")
}

// TestConfigFailure asserts that an error is returned if the config file
// cannot be read.
func TestConfigFailure(t *testing.T) {

	// attempt to read a non-existent config file
	_, err := New(map[string]interface{}{}, "/does-not-exit", false)
	assert.Error(t, err, "failed to error on unreadable config")
}

// TestEmptyEditor asserts that envvars are respected if an editor is not
// specified in the configs
func TestEmptyEditor(t *testing.T) {

	// clear the environment variables
	os.Setenv("VISUAL", "")
	os.Setenv("EDITOR", "")

	// initialize a config
	conf, err := New(map[string]interface{}{}, mock.Path("conf/empty.yml"), false)
	assert.NoError(t, err, "failed to initialize test")

	// set editor, and assert that it is respected
	os.Setenv("EDITOR", "foo")
	conf, err = New(map[string]interface{}{}, mock.Path("conf/empty.yml"), false)
	assert.NoError(t, err, "failed to init configs")
	assert.Equal(t, "foo", conf.Editor, "failed to respect editor")

	// set visual, and assert that it overrides editor
	os.Setenv("VISUAL", "bar")
	conf, err = New(map[string]interface{}{}, mock.Path("conf/empty.yml"), false)
	assert.NoError(t, err, "failed to init configs")
	assert.Equal(t, "bar", conf.Editor, "failed to respect editor")
}

// Asserts that colorization rules are properly respected
func TestColor(t *testing.T) {
	noColorConf := Config{
		Colorize: false,
	}
	assert.Equal(t, false, noColorConf.Colorize)

	colorConf := Config{
		Colorize: true,
	}
	assert.Equal(t, true, colorConf.Colorize)
}

// Asserts that configs are properly initialized
func TestInit(t *testing.T) {
	// Initialize a temporary config file:
	confFile, err := os.CreateTemp("", "cheat-test")
	assert.NoError(t, err, "failed to create temp file")
	defer os.Remove(confFile.Name())

	// Initialize the config file:
	conf := "mock config data"
	err = Init(confFile.Name(), conf)
	assert.NoError(t, err, "failed to init config file")

	// Read back the config file contents:
	bytes, err := os.ReadFile(confFile.Name())
	assert.NoError(t, err, "failed to read config file")

	// Assert that the contents were written correctly:
	got := string(bytes)
	assert.Equal(t, conf, got, "failed to write config")
}

// Asserts that `Path` identifies non-existent config files
func TestPathConfigNotExists(t *testing.T) {
	paths := []string{"/cheat-test-conf-does-not-exist"}
	_, err := Path(paths)
	assert.Error(t, err, "failed to identify non-existent config file")
}

// Asserts that `Path` identifies existent config files
func TestPathConfigExists(t *testing.T) {
	// Initialize a temporary config file:
	confFile, err := os.CreateTemp("", "cheat-test")
	assert.NoError(t, err, "failed to create temp file")
	defer os.Remove(confFile.Name())

	paths := []string{
		"/cheat-test-conf-does-not-exist",
		confFile.Name(),
	}
	got, err := Path(paths)
	assert.NoError(t, err, "failed to identify config file")
	assert.Equal(t, confFile.Name(), got, "failed to return config path")
}

// Asserts that the proper config paths are returned on *nix platforms
func TestValidatePathsNix(t *testing.T) {
	home := "/home/foo"
	envvars := map[string]string{
		"XDG_CONFIG_HOME": "/home/bar",
	}
	oses := []string{
		"android",
		"darwin",
		"freebsd",
		"linux",
	}

	for _, os := range oses {
		paths, err := Paths(os, home, envvars)
		assert.NoError(t, err, "paths returned an error")
		want := []string{
			"/home/bar/cheat/conf.yml",
			"/home/foo/.config/cheat/conf.yml",
			"/home/foo/.cheat/conf.yml",
			"/etc/cheat/conf.yml",
		}
		assert.Equal(t, want, paths, "failed to return exptected paths")
	}
}

// Asserts that the proper config paths are returned on *nix platforms when
// `XDG_CONFIG_HOME is not set
func TestValidatePathsNixNoXDG(t *testing.T) {
	home := "/home/foo"
	envvars := map[string]string{}
	oses := []string{
		"darwin",
		"freebsd",
		"linux",
	}
	for _, os := range oses {
		paths, err := Paths(os, home, envvars)
		assert.NoError(t, err, "paths returned an error")
		want := []string{
			"/home/foo/.config/cheat/conf.yml",
			"/home/foo/.cheat/conf.yml",
			"/etc/cheat/conf.yml",
		}
		assert.Equal(t, want, paths, "failed to return exptected paths")
	}
}

// Asserts that the proper config paths are returned on Windows platforms
func TestValidatePathsWindows(t *testing.T) {
	home := "not-used-on-windows"
	envvars := map[string]string{
		"APPDATA":     "/apps",
		"PROGRAMDATA": "/programs",
	}
	paths, err := Paths("windows", home, envvars)
	assert.NoError(t, err, "paths returned an error")
	want := []string{
		"/apps/cheat/conf.yml",
		"/programs/cheat/conf.yml",
	}
	assert.Equal(t, want, paths, "failed to return exptected paths")
}

// Asserts that an error is returned on unsupported platforms
func TestValidatePathsUnsupported(t *testing.T) {
	_, err := Paths("unsupported", "", map[string]string{})
	assert.Error(t, err, "failed to return error on unsupported platform")
}

// Asserts that the proper config path is returned when `CHEAT_CONFIG_PATH` is
// explicitly specified.
func TestValidatePathsCheatConfigPath(t *testing.T) {
	home := "/home/foo"
	envvars := map[string]string{
		"XDG_CONFIG_HOME":   "/home/bar",
		"CHEAT_CONFIG_PATH": "/home/baz/conf.yml",
	}
	paths, err := Paths("linux", home, envvars)
	assert.NoError(t, err, "paths returned an error")
	want := []string{"/home/baz/conf.yml"}
	assert.Equal(t, want, paths, "failed to return expected paths")
}

// Asserts that valid configs are validated successfully
func TestValidateCorrect(t *testing.T) {
	conf := Config{
		Colorize:  true,
		Editor:    "vim",
		Formatter: "terminal16m",
		Notebooks: []notebook.Notebook{
			{
				Name:     "foo",
				Path:     "/foo",
				ReadOnly: false,
				Tags:     []string{},
			},
		},
	}
	err := conf.Validate()
	assert.NoError(t, err, "failed to validate valid config")
}

// Asserts that configs with unspecified editors are invalidated
func TestInvalidateMissingEditor(t *testing.T) {
	conf := Config{
		Colorize:  true,
		Formatter: "terminal16m",
		Notebooks: []notebook.Notebook{
			{
				Name:     "foo",
				Path:     "/foo",
				ReadOnly: false,
				Tags:     []string{},
			},
		},
	}
	err := conf.Validate()
	assert.Error(t, err, "failed to invalidate config with unspecified editor")
}

// Asserts that configs without notepaths are invalidated
func TestInvalidateMissingCheatpaths(t *testing.T) {
	conf := Config{
		Colorize:  true,
		Editor:    "vim",
		Formatter: "terminal16m",
	}
	err := conf.Validate()
	assert.Error(t, err, "failed to invalidate config without cheatpaths")
}

// Asserts that configs which contain invalid formatters are invalidated
func TestMissingInvalidFormatters(t *testing.T) {
	conf := Config{
		Colorize: true,
		Editor:   "vim",
	}
	err := conf.Validate()
	assert.Error(t, err, "failed to invalidate config without formatter")
}

// Asserts that configs which contain notepaths with duplcated names are
// invalidated
func TestInvalidateDuplicateCheatpathNames(t *testing.T) {
	conf := Config{
		Colorize:  true,
		Editor:    "vim",
		Formatter: "terminal16m",
		Notebooks: []notebook.Notebook{
			{
				Name:     "foo",
				Path:     "/foo",
				ReadOnly: false,
				Tags:     []string{},
			},
			{
				Name:     "foo",
				Path:     "/bar",
				ReadOnly: false,
				Tags:     []string{},
			},
		},
	}
	err := conf.Validate()
	assert.Error(t, err, "failed to invalidate config with cheatpaths with duplicate names")
}

// Asserts that configs which contain notepaths with duplcated paths are
// invalidated
func TestInvalidateDuplicateCheatpathPaths(t *testing.T) {
	conf := Config{
		Colorize:  true,
		Editor:    "vim",
		Formatter: "terminal16m",
		Notebooks: []notebook.Notebook{
			{
				Name:     "foo",
				Path:     "/foo",
				ReadOnly: false,
				Tags:     []string{},
			},
			{
				Name:     "bar",
				Path:     "/foo",
				ReadOnly: false,
				Tags:     []string{},
			},
		},
	}
	err := conf.Validate()
	assert.Error(t, err, "failed to invalidate config with cheatpaths with duplicate paths")
}
