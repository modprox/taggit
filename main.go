// Author hoenig

// Command taggit provides a convenience wrapper around `git tag` commands.
package main

import (
	"fmt"
	"os"
	"strings"
)

func usage(code int) {
	fmt.Fprintf(os.Stderr, "usage: taggit [list, zero, patch, minor, major]\n")
	os.Exit(code)
}

func main() {
	if len(os.Args) != 2 {
		usage(1)
	}

	command := os.Args[1]

	tags := listTags()

	var err error
	switch command {
	case "help":
		usage(0)
	case "list":
		list(tags)
	case "zero":
		err = zero(tags)
	case "patch":
		err = patch(tags)
	case "minor":
		err = minor(tags)
	case "major":
		err = major(tags)

	default:
		usage(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed: %v\n", err)
		os.Exit(1)
	}
}

func list(tags []Tag) {
	var b strings.Builder
	for _, tag := range tags {
		b.WriteString(tag.String())
		b.WriteString("\n")
	}
	fmt.Print(b.String())
}

func zero(tags []Tag) error {
	if len(tags) > 0 {
		return fmt.Errorf("refusing to generate zero tag (v0.0.0) when other semver tags already exist")
	}
	createTag(ZeroTag)

	return nil
}

func patch(tags []Tag) error {
	if len(tags) < 1 {
		return fmt.Errorf("refusing to bump patch with no pre-existing tag")
	}

	latest := tags[0]

	createTag(Tag{
		Major: latest.Major,
		Minor: latest.Minor,
		Patch: latest.Patch + 1,
	})

	return nil
}

func minor(tags []Tag) error {
	if len(tags) < 1 {
		return fmt.Errorf("refusing to bump minor with no pre-existing tag")
	}

	latest := tags[0]

	createTag(Tag{
		Major: latest.Major,
		Minor: latest.Minor + 1,
		Patch: 0,
	})

	return nil
}

func major(tags []Tag) error {
	if len(tags) < 1 {
		return fmt.Errorf("refusing to bump major with no pre-existnig tag")
	}

	latest := tags[0]

	createTag(Tag{
		Major: latest.Major + 1,
		Minor: 0,
		Patch: 0,
	})

	return nil
}
