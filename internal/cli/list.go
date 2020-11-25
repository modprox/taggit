package cli

import (
	git5 "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"gophers.dev/pkgs/semantic"

	"oss.indeed.com/go/taggit/internal/tags"
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i TagLister -s _mock.go

// A TagLister is used to list tags associated with the git repository of a
// Go module.
type TagLister interface {
	// ListRepoTags lists the tags on a given repository.
	ListRepoTags() (tags.Taxonomy, error)
}

func NewTagLister(r *git5.Repository) TagLister {
	return &tagLister{
		repository: r,
	}
}

type tagLister struct {
	repository *git5.Repository
}

func (tl *tagLister) ListRepoTags() (tags.Taxonomy, error) {
	iter, err := tl.repository.Tags()
	if err != nil {
		return nil, err
	}

	groups := make(tags.Taxonomy)
	if err := iter.ForEach(func(ref *plumbing.Reference) error {
		if tag, ok := semantic.Parse(ref.Name().Short()); ok {
			groups.Add(tag)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	groups.Sort()
	return groups, nil
}
