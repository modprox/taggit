package git

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/modprox/taggit/tags"
	"github.com/pkg/errors"
)

// newest version to oldest version
func ListTags(cmd Cmd) ([]tags.Tag, error) {
	output, err := cmd.Run([]string{"tag", "-l"}, 10*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list tags")
	}

	lines := strings.Split(output, "\n")

	parsedTags := make([]tags.Tag, 0, len(lines))
	for _, line := range lines {
		if tag, ok := tags.Parse(line); ok {
			parsedTags = append(parsedTags, tag)
		}
	}

	sortable := tags.BySemver(parsedTags)
	sort.Sort(sort.Reverse(sortable))

	return parsedTags, nil
}

func CreateTag(cmd Cmd, tag tags.Tag) error {
	fmt.Println("taggit: creating tag:", tag)

	if _, err := cmd.Run([]string{"tag", tag.String()}, 3*time.Second); err != nil {
		return err
	}

	if _, err := cmd.Run([]string{"push", "origin", tag.String()}, 1*time.Minute); err != nil {
		return err
	}

	return nil
}
