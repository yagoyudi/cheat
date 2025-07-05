// Implements functions pertaining to loading, sorting, filtering, and tagging
// notes.
package notes

import (
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/yagoyudi/note/internal/note"
	"github.com/yagoyudi/note/internal/notebook"
	"github.com/yagoyudi/note/internal/repo"
)

// Applies notes "overrides", resolving title conflicts that exist among
// notebooks by preferring more local notes over less local notes
func Consolidate(notebooks []map[string]note.Note) map[string]note.Note {
	consolidated := make(map[string]note.Note)
	for _, notebook := range notebooks {
		maps.Copy(consolidated, notebook)
	}
	return consolidated
}

// Organizes the notes into an alphabetically-sorted slice
func Sort(notes map[string]note.Note) []note.Note {
	var titles []string
	for title := range notes {
		titles = append(titles, title)
	}

	sort.Strings(titles)
	sorted := []note.Note{}
	for _, title := range titles {
		sorted = append(sorted, notes[title])
	}
	return sorted
}

// Returns a slice of all tags in use in any note
func Tags(notebooks []map[string]note.Note) []string {
	// Create a map of all tags in use in any sheet:
	tags := make(map[string]bool)

	// Iterate over all tags on all sheets on all cheatpaths:
	for _, notebook := range notebooks {
		for _, note := range notebook {
			for _, tag := range note.Tags {
				tags[tag] = true
			}
		}
	}

	sorted := []string{}
	for tag := range tags {
		sorted = append(sorted, tag)
	}
	slices.Sort(sorted)
	return sorted
}

// Filter filters notes that do not match `tag(s)`
func Filter(notebooks []map[string]note.Note, tags []string) []map[string]note.Note {
	// Buffer a map of filtered notes:
	filtered := make([]map[string]note.Note, 0, len(notebooks))
	for _, notes := range notebooks {

		// Create a map of notes for each notepath:
		pathFiltered := make(map[string]note.Note)

		for title, note := range notes {
			// Assume that the note should be kept (ie, should not be filtered):
			keep := true

			// Iterate over each tag. If the note does not match *all* tags,
			// filter it out:
			for _, tag := range tags {
				if !note.TaggedWith(strings.TrimSpace(tag)) {
					keep = false
				}
			}

			// If the note does match all tags, it passes the filter:
			if keep {
				pathFiltered[title] = note
			}
		}
		filtered = append(filtered, pathFiltered)
	}
	return filtered
}

// Produces a map of note names to filesystem paths
func Load(notebooks []notebook.Notebook) ([]map[string]note.Note, error) {
	notesByName := make([]map[string]note.Note, len(notebooks))
	for _, notebook := range notebooks {
		noteByName := make(map[string]note.Note)

		// Recursively iterate over the notebook, and load each note
		// encountered along the way:
		err := filepath.Walk(
			notebook.Path,
			func(path string, info os.FileInfo, err error) error {
				// Fail if an error occurred while walking the directory:
				if err != nil {
					return fmt.Errorf("failed to walk path: %v", err)
				}

				// Don't register directories as notes:
				if info.IsDir() {
					return nil
				}

				// Calculate the notes's "name" (the phrase with which it may be
				// accessed. Eg: `cheat tar` - `tar` is the name):
				name := strings.TrimPrefix(
					strings.TrimPrefix(path, notebook.Path),
					string(os.PathSeparator),
				)

				// Don't walk the `.git` directory:
				skip, err := repo.GitDir(path)
				if err != nil {
					return fmt.Errorf("failed to identify .git directory: %v", err)
				}
				if skip {
					return fs.SkipDir
				}

				note, err := note.New(name, notebook.Name, path, notebook.Tags, notebook.ReadOnly)
				if err != nil {
					return fmt.Errorf("failed to load note: %s, path: %s, err: %v", name, path, err)
				}
				noteByName[name] = note
				return nil
			},
		)
		if err != nil {
			return notesByName, fmt.Errorf("failed to load notes: %v", err)
		}
		notesByName = append(notesByName, noteByName)
	}
	return notesByName, nil
}
