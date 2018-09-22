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

	gitCmd := git.New("git")

	tags, err := git.ListTags(gitCmd)
	if err != nil {
		die(err)
	}

	switch command {
	case "help":
		cli.Usage(0)
	case "list":
		cli.List(tags)
	case "zero":
		err = cli.Zero(gitCmd, tags)
	case "patch":
		err = cli.Patch(gitCmd, tags)
	case "minor":
		err = cli.Minor(gitCmd, tags)
	case "major":
		err = cli.Major(gitCmd, tags)

	default:
		cli.Usage(1)
	}

	if err != nil {
		die(err)
	}
}

func die(err error) {
	fmt.Fprintf(os.Stderr, "failed: %v\n", err)
	os.Exit(1)
}
