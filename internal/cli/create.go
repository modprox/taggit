package cli

import (
	git5 "github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"gophers.dev/pkgs/semantic"
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i TagCreator -s _mock.go

// A TagCreator is used to create git tags.
type TagCreator interface {
	// CreateTag creates the given tag in the repository.
	CreateTag(semantic.Tag) error
}

type tagCreator struct {
	repository *git5.Repository
}

func NewTagCreator(r *git5.Repository) TagCreator {
	return &tagCreator{
		repository: r,
	}
}

func (tc *tagCreator) CreateTag(tag semantic.Tag) error {
	head, err := tc.repository.Head()
	if err != nil {
		return errors.Wrap(err, "could not find head")
	}

	if _, err := tc.repository.CreateTag(
		tag.String(),
		head.Hash(),
		nil, // options, maybe include a message
	); err != nil {
		return errors.Wrap(err, "could not create tag")
	}

	return nil
}
