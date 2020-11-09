package cli

import (
	"context"
	"fmt"
	"time"

	git5 "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/pkg/errors"
	"gophers.dev/pkgs/semantic"
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i TagPusher -s _mock.go

// A TagPusher is used to push tags to a remote.
type TagPusher interface {
	// PushTag pushes the given tag to the default remote.
	PushTag(semantic.Tag) error
}

type tagPusher struct {
	repository *git5.Repository
}

func NewTagPusher(r *git5.Repository) TagPusher {
	return &tagPusher{
		repository: r,
	}
}

func (tp *tagPusher) PushTag(tag semantic.Tag) error {
	remotes, err := tp.repository.Remotes()
	if err != nil {
		return errors.Wrap(err, "could not find remotes")
	}

	if len(remotes) == 0 {
		return nil
	}

	label := tag.String()
	refSpec := fmt.Sprintf("refs/tags/%s:refs/tags/%s", label, label)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := tp.repository.PushContext(ctx, &git5.PushOptions{
		RemoteName: "", // default
		RefSpecs:   []config.RefSpec{config.RefSpec(refSpec)},
	}); err != nil {
		return errors.Wrap(err, "could not push tag")
	}

	return nil
}
