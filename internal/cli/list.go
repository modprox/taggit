package cli

import (
	"bufio"
	"strings"
	"time"

	"github.com/pkg/errors"

	"gophers.dev/pkgs/semantic"

	"oss.indeed.com/go/taggit/internal/git"
	"oss.indeed.com/go/taggit/internal/tags"
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i TagLister -s _mock.go

// A TagLister is used to list tags associated with the git repository of a
// Go module.
type TagLister interface {
	// ListRepoTags lists the tags on a given repository.
	ListRepoTags() (tags.Taxonomy, error)
}

func NewTagLister(gitCmd git.Cmd) TagLister {
	return &tagLister{
		timeout: 3 * time.Second,
		gitCmd:  gitCmd,
	}
}

type tagLister struct {
	timeout time.Duration
	gitCmd  git.Cmd
}

var repoTagsArgs = []string{"tag", "--list"}

func (tl *tagLister) ListRepoTags() (tags.Taxonomy, error) {
	output, err := tl.gitCmd.Run(repoTagsArgs, tl.timeout)
	if err != nil {
		return nil, errors.Wrap(err, "unable to list repo tags")
	}

	groups, err := parseTagList(output)
	if err != nil {
		return nil, errors.Wrap(err, "parse tags list")
	}

	return groups, nil
}

// Parse the output of `git tag --list` into a set of tags grouped
// such that all extensions are associated with a core version.
func parseTagList(s string) (tags.Taxonomy, error) {
	groups := make(tags.Taxonomy)
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if tag, ok := semantic.Parse(line); ok {
			groups.Add(tag)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	groups.Sort()
	return groups, nil
}
