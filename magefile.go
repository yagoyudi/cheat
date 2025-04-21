//go:build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
)

var Default = Build

func Build() error {
	err := sh.RunV("go", "build", "-o", "bin/cheat", "./cmd/cheat")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Println("Check the ./bin folder!")
	return nil
}

func Test() error {
	return sh.RunV("go", "test", "-v", "./...")
}

func Clean() error {
	return sh.RunV("rm", "-rf", "bin")
}
