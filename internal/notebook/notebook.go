// Implements functions pertaining to note file path management.
package notebook

import "fmt"

// Encapsulates notebook information
type Notebook struct {
	Name     string   `yaml:"name"`
	Path     string   `yaml:"path"`
	ReadOnly bool     `yaml:"readonly"`
	Tags     []string `yaml:"tags"`
}

// Returns an error if the notepath is invalid
func (n *Notebook) Validate() error {
	if n.Name == "" {
		return fmt.Errorf("cheatpath: invalid cheatpath: name must be specified")
	}
	if n.Path == "" {
		return fmt.Errorf("cheatpath: invalid cheatpath: path must be specified")
	}

	return nil
}

// Filters all notepaths that are not named `name`
func Filter(paths []Notebook, name string) ([]Notebook, error) {
	for _, path := range paths {
		if path.Name == name {
			return []Notebook{path}, nil
		}
	}
	return []Notebook{}, fmt.Errorf("cheatpath: cheatpath does not exist: %s", name)
}

// Returns a writeable NotePath
func Writeable(notebooks []Notebook) (Notebook, error) {
	// Iterate backwards over the notebooks:
	// NOTE: we're going backwards because we assume that the most "local"
	// notepath will be specified last in the configs.
	for i := len(notebooks) - 1; i >= 0; i-- {
		if !notebooks[i].ReadOnly {
			return notebooks[i], nil
		}
	}
	return Notebook{}, fmt.Errorf("no writeable cheatpaths found")
}
