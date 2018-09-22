// Command taggit provides a convenience wrapper around `git tag` commands.
package main

import (
	"fmt"
	"os"

	"github.com/modprox/taggit/internal/cli"
	"github.com/modprox/taggit/internal/git"
)

func main() {
	if len(os.Args) != 2 {
		cli.Usage(1)
	}

	command := os.Args[1]

	tags := git.ListTags()

	var err error
	switch command {
	case "help":
		cli.Usage(0)
	case "list":
		cli.List(tags)
	case "zero":
		err = cli.Zero(tags)
	case "patch":
		err = cli.Patch(tags)
	case "minor":
		err = cli.Minor(tags)
	case "major":
		err = cli.Major(tags)

	default:
		cli.Usage(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed: %v\n", err)
		os.Exit(1)
	}
}
