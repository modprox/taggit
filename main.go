// Command taggit provides a convenience wrapper around `git tag` commands.
package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/modprox/taggit/internal/cli"
	"github.com/modprox/taggit/internal/git"
	"github.com/modprox/taggit/internal/publish"
)

func main() {
	if len(os.Args) != 2 {
		cli.Usage(1)
	}

	command := os.Args[1]

	gitCmd := git.New("git")
	var output bytes.Buffer

	// modFinder := publish.NewModFinder()
	// pubber := publish.New(publish.Discard(), modFinder)
	publisher := publish.Discard()

	tool := cli.NewTool(&output, gitCmd, publisher)

	tags, err := git.ListTags(gitCmd)
	if err != nil {
		die(err)
	}

	switch command {
	case "help":
		cli.Usage(0)
	case "list":
		err = tool.List(tags)
	case "zero":
		err = tool.Zero(tags)
	case "patch":
		err = tool.Patch(tags)
	case "minor":
		err = tool.Minor(tags)
	case "major":
		err = tool.Major(tags)
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
