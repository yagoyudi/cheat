// Implements mock functions used in unit-tests.
package mock

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"
)

// Returns the absolute path to the specified mock file.
func Path(filename string) string {
	// Determine the path of this file during runtime:
	_, thisfile, _, _ := runtime.Caller(0)
	file, err := filepath.Abs(path.Join(filepath.Dir(thisfile), "../mocks", filename))
	if err != nil {
		panic(fmt.Errorf("failed to resolve mock path: %v", err))
	}
	return file
}
