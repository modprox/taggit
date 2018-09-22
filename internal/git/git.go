package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/modprox/taggit/internal/tags"
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
func ListTags() []tags.Tag {
	output := git([]string{"tag", "-l"}, 10*time.Second)
	lines := strings.Split(output, "\n")

	parsedTags := make([]tags.Tag, 0, len(lines))
	for _, line := range lines {
		if tag, ok := tags.Parse(line); ok {
			parsedTags = append(parsedTags, tag)
		}
	}

	sortable := tags.BySemver(parsedTags)
	sort.Sort(sort.Reverse(sortable))

	return parsedTags
}

func CreateTag(tag tags.Tag) {
	fmt.Println("taggit: creating tag:", tag)

	_ = git([]string{"tag", tag.String()}, 3*time.Second)
	_ = git([]string{"push", "origin", tag.String()}, 1*time.Minute)
}
