package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yagoyudi/cheat/internal/cheatpath"
	"github.com/yagoyudi/cheat/internal/mock"
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
	want := []cheatpath.Cheatpath{
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
	assert.Equal(t, want, conf.Cheatpaths, "failed to return expected results")
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
