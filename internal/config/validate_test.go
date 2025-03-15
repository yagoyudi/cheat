package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yagoyudi/cheat/internal/cheatpath"
)

// TestValidateCorrect asserts that valid configs are validated successfully
func TestValidateCorrect(t *testing.T) {

	// mock a config
	conf := Config{
		Colorize:  true,
		Editor:    "vim",
		Formatter: "terminal16m",
		Cheatpaths: []cheatpath.Cheatpath{
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

// TestInvalidateMissingEditor asserts that configs with unspecified editors
// are invalidated
func TestInvalidateMissingEditor(t *testing.T) {

	// mock a config
	conf := Config{
		Colorize:  true,
		Formatter: "terminal16m",
		Cheatpaths: []cheatpath.Cheatpath{
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

// TestInvalidateMissingCheatpaths asserts that configs without cheatpaths are
// invalidated
func TestInvalidateMissingCheatpaths(t *testing.T) {

	// mock a config
	conf := Config{
		Colorize:  true,
		Editor:    "vim",
		Formatter: "terminal16m",
	}

	err := conf.Validate()
	assert.Error(t, err, "failed to invalidate config without cheatpaths")
}

// TestMissingInvalidFormatters asserts that configs which contain invalid
// formatters are invalidated
func TestMissingInvalidFormatters(t *testing.T) {

	// mock a config
	conf := Config{
		Colorize: true,
		Editor:   "vim",
	}

	err := conf.Validate()
	assert.Error(t, err, "failed to invalidate config without formatter")
}

// TestInvalidateDuplicateCheatpathNames asserts that configs which contain
// cheatpaths with duplcated names are invalidated
func TestInvalidateDuplicateCheatpathNames(t *testing.T) {

	// mock a config
	conf := Config{
		Colorize:  true,
		Editor:    "vim",
		Formatter: "terminal16m",
		Cheatpaths: []cheatpath.Cheatpath{
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

// TestInvalidateDuplicateCheatpathPaths asserts that configs which contain
// cheatpaths with duplcated paths are invalidated
func TestInvalidateDuplicateCheatpathPaths(t *testing.T) {

	// mock a config
	conf := Config{
		Colorize:  true,
		Editor:    "vim",
		Formatter: "terminal16m",
		Cheatpaths: []cheatpath.Cheatpath{
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
