package cli

import (
	"time"

	"github.com/pkg/errors"

	"gophers.dev/pkgs/semantic"

	"oss.indeed.com/go/taggit/internal/git"
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i TagCreator -s _mock.go

type TagCreator interface {
	CreateTag(semantic.Tag) error
}

type tagCreator struct {
	timeout time.Duration
	gitCmd  git.Cmd
}

func NewTagCreator(gitCmd git.Cmd) TagCreator {
	return &tagCreator{
		timeout: 3 * time.Second,
		gitCmd:  gitCmd,
	}
}

func (tc *tagCreator) CreateTag(tag semantic.Tag) error {
	createArgs := createTagArgs(tag)
	if _, err := tc.gitCmd.Run(createArgs, tc.timeout); err != nil {
		return errors.Wrap(err, "unable to create tag")
	}

	// for now, just assume there is an upstream to push to, since that
	// will be the case like 99.99% of the time
	pushArgs := pushTagArgs(tag)
	if _, err := tc.gitCmd.Run(pushArgs, tc.timeout); err != nil {
		return errors.Wrap(err, "unable to push tag")
	}

	return nil
}

func createTagArgs(tag semantic.Tag) []string {
	return []string{"tag", tag.String()}
}

func pushTagArgs(tag semantic.Tag) []string {
	return []string{"push", "origin", tag.String()}
}
