// Implements functions pertaining to parsing, searching, and displaying notes.
package note

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/yagoyudi/note/internal/config"
	"gopkg.in/yaml.v3"
)

// Encapsulates cheatsheet header data
type noteHeader struct {
	Tags   []string
	Syntax string
}

// Encapsulates note information
type Note struct {
	Name     string
	Notebook string
	Path     string
	Body     string
	Tags     []string
	Syntax   string
	ReadOnly bool
}

// Initializes a new note
func New(name string, notebook string, path string, tags []string, readOnly bool) (Note, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Note{}, fmt.Errorf("failed to read file: %s, %v", path, err)
	}

	header, body, err := parse(string(raw))
	if err != nil {
		return Note{}, fmt.Errorf("failed to parse header: %v", err)
	}

	// Merge the sheet-specific tags into the cheatpath tags:
	tags = append(tags, header.Tags...)

	// Sort strings so they pretty-print nicely:
	sort.Strings(tags)

	return Note{
		Name:     name,
		Notebook: notebook,
		Path:     path,
		Body:     body,
		Tags:     tags,
		Syntax:   header.Syntax,
		ReadOnly: readOnly,
	}, nil
}

// Parses note header
func parse(raw string) (noteHeader, string, error) {
	delim := "---\n"

	// If the raw note does not contain header, pass it through unmodified:
	if !strings.HasPrefix(raw, delim) {
		return noteHeader{}, raw, nil
	}

	// Split the header and body:
	parts := strings.SplitN(raw, delim, 3)

	// Return an error if the header parses into the wrong number of parts:
	if len(parts) != 3 {
		return noteHeader{}, "", fmt.Errorf("failed to delimit header")
	}

	var header noteHeader
	if err := yaml.Unmarshal([]byte(parts[1]), &header); err != nil {
		return noteHeader{}, "", fmt.Errorf("failed to unmarshal header: %v", err)
	}

	return header, parts[2], nil
}

// Returns lines within a note's body that match the regex
func (n *Note) Search(reg *regexp.Regexp) string {
	matches := ""
	for _, line := range strings.Split(n.Body, "\n\n") {
		// exit early if the line doesn't match the regex
		if reg.MatchString(line) {
			matches += line + "\n\n"
		}
	}
	return strings.TrimSpace(matches)
}

// Copies a note to a new location
func (n *Note) Copy(dest string) error {
	// NB: while the `infile` has already been loaded and parsed into a `sheet`
	// struct, we're going to read it again here. This is a bit wasteful, but
	// necessary if we want the "raw" file contents (including the front-matter).
	// This is because the frontmatter is parsed and then discarded when the file
	// is loaded via `sheets.Load`.
	infile, err := os.Open(n.Path)
	if err != nil {
		return fmt.Errorf("failed to open cheatsheet: %s, %v", n.Path, err)
	}
	defer infile.Close()

	// Create any necessary subdirectories:
	dirs := filepath.Dir(dest)
	if dirs != "." {
		if err := os.MkdirAll(dirs, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %s, %v", dirs, err)
		}
	}

	// Create the outfile:
	outfile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create outfile: %s, %v", dest, err)
	}
	defer outfile.Close()

	// Copy file contents:
	_, err = io.Copy(outfile, infile)
	if err != nil {
		return fmt.Errorf("failed to copy file: infile: %s, outfile: %s, err: %v", n.Path, dest, err)
	}

	return nil
}

// Applies syntax-highlighting to a notes's body
func (n *Note) Colorize(conf config.Config) {
	// If the syntax was not specified, default to bash:
	lex := n.Syntax
	if lex == "" {
		lex = "bash"
	}

	// write colorized text into a buffer
	var buf bytes.Buffer
	err := quick.Highlight(&buf, n.Body, lex, conf.Formatter, conf.Style)
	if err != nil {
		// if colorization somehow failed, do nothing
		return
	}

	// Swap the note's body with its colorized equivalent:
	n.Body = buf.String()
}

// Returns true if a sheet was tagged with `target`
func (s *Note) Tagged(target string) bool {
	for _, tag := range s.Tags {
		if tag == target {
			return true
		}
	}
	return false
}
