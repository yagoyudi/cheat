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
		return fmt.Errorf("notebook name must be specified")
	}
	if n.Path == "" {
		return fmt.Errorf("notebook path must be specified")
	}

	return nil
}

// Filters all notebooks that are not named `name`
func Filter(notebooks []Notebook, name string) ([]Notebook, error) {
	for _, notebook := range notebooks {
		if notebook.Name == name {
			return []Notebook{notebook}, nil
		}
	}
	return []Notebook{}, fmt.Errorf("notebook %s does not exist", name)
}

// Returns a writeable notebook
func Writeable(notebooks []Notebook) (Notebook, error) {
	// NOTE: we're going backwards because we assume that the most "local"
	// notebook will be specified last in the configs.
	for i := len(notebooks) - 1; i >= 0; i-- {
		if !notebooks[i].ReadOnly {
			return notebooks[i], nil
		}
	}
	return Notebook{}, fmt.Errorf("no writeable notebook found")
}
