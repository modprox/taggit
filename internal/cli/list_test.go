package cli

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gophers.dev/pkgs/semantic"
	"oss.indeed.com/go/taggit/internal/tags"
)

func Test_TagLister_ListRepoTags_normal(t *testing.T) {
	repo := CreateT(t, []string{
		"v0.0.1",
		"v0.1.0-rc1",
		"v0.1.0-alpha1+bm1",
		"v0.1.0",
		"v1.0.0",
		"v1.0.0-rc1",
	})
	defer CleanupT(t, repo)

	lister := NewTagLister(repo)
	tax, err := lister.ListRepoTags()
	require.NoError(t, err)
	require.Equal(t, tags.Taxonomy{
		tags.NewTriple(0, 0, 1): []semantic.Tag{
			semantic.New(0, 0, 1),
		},
		tags.NewTriple(0, 1, 0): []semantic.Tag{
			semantic.New(0, 1, 0),
			semantic.New2(0, 1, 0, "rc1"),
			semantic.New3(0, 1, 0, "alpha1", "bm1"),
		},
		tags.NewTriple(1, 0, 0): []semantic.Tag{
			semantic.New(1, 0, 0),
			semantic.New2(1, 0, 0, "rc1"),
		},
	}, tax)
}
