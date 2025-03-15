package sheet

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCopyFlat asserts that Copy correctly copies files at a single level of
// depth
func TestCopyFlat(t *testing.T) {

	// mock a cheatsheet file
	text := "this is the cheatsheet text"
	src, err := os.CreateTemp("", "foo-src")
	assert.NoError(t, err, "failed to mock cheatsheet")
	defer src.Close()
	defer os.Remove(src.Name())

	_, err = src.WriteString(text)
	assert.NoError(t, err, "failed to write to mock cheatsheet")

	// mock a cheatsheet struct
	sheet, err := New("foo", "community", src.Name(), []string{}, false)
	assert.NoError(t, err, "failed to init cheatsheet")

	// compute the outfile's path
	outpath := path.Join(os.TempDir(), sheet.Title)
	defer os.Remove(outpath)

	// attempt to copy the cheatsheet
	err = sheet.Copy(outpath)
	assert.NoError(t, err, "failed to copy cheatsheet")

	// assert that the destination file contains the correct text
	got, err := os.ReadFile(outpath)
	assert.NoError(t, err, "failed to read destination file")
	assert.Equal(t, text, string(got), "destination file contained wrong text")
}

// TestCopyDeep asserts that Copy correctly copies files at several levels of
// depth
func TestCopyDeep(t *testing.T) {

	// mock a cheatsheet file
	text := "this is the cheatsheet text"
	src, err := os.CreateTemp("", "foo-src")
	assert.NoError(t, err, "failed to mock cheatsheet")
	defer src.Close()
	defer os.Remove(src.Name())

	_, err = src.WriteString(text)
	assert.NoError(t, err, "failed to write to mock cheatsheet")

	// mock a cheatsheet struct
	sheet, err := New(
		"/cheat-tests/alpha/bravo/foo",
		"community",
		src.Name(),
		[]string{},
		false,
	)
	assert.NoError(t, err, "failed to init cheatsheet")

	// compute the outfile's path
	outpath := path.Join(os.TempDir(), sheet.Title)
	defer os.RemoveAll(path.Join(os.TempDir(), "cheat-tests"))

	// attempt to copy the cheatsheet
	err = sheet.Copy(outpath)
	assert.NoError(t, err, "failed to copy cheatsheet")

	// assert that the destination file contains the correct text
	got, err := os.ReadFile(outpath)
	assert.NoError(t, err, "failed to read destination file")
	assert.Equal(t, text, string(got), "destination file contained wrong text")
}
