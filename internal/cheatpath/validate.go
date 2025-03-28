package cheatpath

import (
	"fmt"
)

// Validate returns an error if the cheatpath is invalid
func (c *Cheatpath) Validate() error {

	if c.Name == "" {
		return fmt.Errorf("cheatpath: invalid cheatpath: name must be specified")
	}
	if c.Path == "" {
		return fmt.Errorf("cheatpath: invalid cheatpath: path must be specified")
	}

	return nil
}
