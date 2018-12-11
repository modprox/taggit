// Command taggit provides a convenience wrapper around `git tag` commands.
package main

import (
	"fmt"
	"os"

	"github.com/modprox/taggit/internal/cli"
)

func main() {
	starter, err := cli.New(os.Args)
	if err != nil {
		fmt.Printf("failed to initialize with error: %v\n", err)
		os.Exit(1)
	}

	if err := starter.Start(); err != nil {
		fmt.Printf("failed to execute with error: %v\n", err)
		os.Exit(1)
	}
}
