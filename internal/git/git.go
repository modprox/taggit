package git

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"gophers.dev/pkgs/semantic"
)

var currentBranchArgs = []string{"rev-parse", "--abbrev-ref", "HEAD"}
var listTagsBase = []string{"tag", "--merged"}

func listTagsArgs(branch string) []string {
	complete := append(listTagsBase, branch)
	return complete
}

func getBranch(cmd Cmd) (string, error) {
	output, err := cmd.Run(currentBranchArgs, 10*time.Second)
	if err != nil {
		return "", errors.Wrap(err, "unable to get current branch")
	}
	return strings.TrimSpace(output), nil
}

// List tags executes git to list tags on the current branch.
//
// It is important to list tags only on the current branch in order to accommodate
// repositories which use branches to distinguish between module major versions.
func ListTags(cmd Cmd) (string, error) {
	branch, err := getBranch(cmd)
	if err != nil {
		return "", errors.Wrap(err, "list tags failed")
	}

	args := listTagsArgs(branch)
	output, err := cmd.Run(args, 10*time.Second)
	if err != nil {
		return "", errors.Wrap(err, "list tags failed")
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
