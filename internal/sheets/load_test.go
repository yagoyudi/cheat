package sheets

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yagoyudi/cheat/internal/cheatpath"
	"github.com/yagoyudi/cheat/internal/mock"
)

// TestLoad asserts that sheets on valid cheatpaths can be loaded successfully
func TestLoad(t *testing.T) {

	// mock cheatpaths
	cheatpaths := []cheatpath.Cheatpath{
		{
			Name:     "community",
			Path:     path.Join(mock.Path("cheatsheets"), "community"),
			ReadOnly: true,
		},
		{
			Name:     "personal",
			Path:     path.Join(mock.Path("cheatsheets"), "personal"),
			ReadOnly: false,
		},
	}

	// load cheatsheets
	sheets, err := Load(cheatpaths)
	assert.NoError(t, err, "failed to load cheatsheets")

	// assert that the correct number of sheets loaded
	// (sheet load details are tested in `sheet_test.go`)
	want := 4
	assert.Equal(t, want, len(sheets), "failed to load correct number of cheatsheets")
}

// TestLoadBadPath asserts that an error is returned if a cheatpath is invalid
func TestLoadBadPath(t *testing.T) {

	// mock a bad cheatpath
	cheatpaths := []cheatpath.Cheatpath{
		{
			Name:     "badpath",
			Path:     "/cheat/test/path/does/not/exist",
			ReadOnly: true,
		},
	}

	// attempt to load the cheatpath
	_, err := Load(cheatpaths)
	assert.Error(t, err, "failed to reject invalid cheatpath")
}
