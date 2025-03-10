// Package main serves as the executable entrypoint.
package main

import (
	"fmt"
	"os"

	"github.com/yagoyudi/cheat/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "cheat: %v", err)
		os.Exit(1)
	}
}
