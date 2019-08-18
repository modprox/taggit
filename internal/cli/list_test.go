package cli

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"gophers.dev/pkgs/semantic"

	"oss.indeed.com/go/taggit/internal/git"
	"oss.indeed.com/go/taggit/internal/tags"
)

const (
	someTags = `
v0.0.1
v0.1.0-rc1
v0.1.0-alpha1+bm1
v0.1.0
v1.0.0
v1.0.0-rc1
`
)

func Test_TagLister_ListRepoTags_normal(t *testing.T) {
	gitCmd := git.NewCmdMock(t)
	defer gitCmd.MinimockFinish()

	gitCmd.RunMock.Expect(
		[]string{"tag", "--list"},
		3*time.Second,
	).Return(someTags, nil)

	tl := NewTagLister(gitCmd)

	tax, err := tl.ListRepoTags()
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

func Test_TagLister_ListRepoTags_errList(t *testing.T) {
	gitCmd := git.NewCmdMock(t)
	defer gitCmd.MinimockFinish()

	gitCmd.RunMock.Expect(
		[]string{"tag", "--list"},
		3*time.Second,
	).Return("", errors.New("git list error"))

	tl := NewTagLister(gitCmd)

	_, err := tl.ListRepoTags()
	require.EqualError(t, err, "unable to list repo tags: git list error")
}
