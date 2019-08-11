package git

import (
	"fmt"
	"time"

	"github.com/pkg/errors"

	"gophers.dev/pkgs/semantic"
)

// newest version to oldest version
func ListTags(cmd Cmd) (string, error) {
	output, err := cmd.Run([]string{"tag", "-l"}, 10*time.Second)
	if err != nil {
		return "", errors.Wrap(err, "git failed to list tags")
	}
	return output, nil
}

func CreateTag(cmd Cmd, tag semantic.Tag) error {
	fmt.Println("taggit: creating tag:", tag)

	if _, err := cmd.Run([]string{"tag", tag.String()}, 3*time.Second); err != nil {
		return err
	}

	if _, err := cmd.Run([]string{"push", "origin", tag.String()}, 1*time.Minute); err != nil {
		return err
	}

	return nil
}
