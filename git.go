package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// git tag -l
func git(args []string, timeout time.Duration) string {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Env = os.Environ() // use the tty's environment
	bs, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "git command failed: %v\n", err)
		fmt.Fprintf(os.Stderr, "output: %s\n", string(bs))
		os.Exit(1)
	}
	return string(bs)
}

// newest version to oldest version
func listTags() []Tag {
	output := git([]string{"tag", "-l"}, 10*time.Second)
	lines := strings.Split(output, "\n")

	tags := make([]Tag, 0, len(lines))
	for _, line := range lines {
		if tag, ok := Parse(line); ok {
			tags = append(tags, tag)
		}
	}

	sortable := TagsBySemver(tags)
	sort.Sort(sort.Reverse(sortable))

	return tags
}

func createTag(tag Tag) {
	fmt.Println("taggit: creating tag:", tag)

	_ = git([]string{"tag", tag.String()}, 3*time.Second)
	_ = git([]string{"push", "origin", tag.String()}, 1*time.Minute)
}
