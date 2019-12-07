package cli

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"gophers.dev/pkgs/semantic"

	"oss.indeed.com/go/taggit/internal/git"
)

func Test_TagCreator_CreateTag(t *testing.T) {
	gitCmd := git.NewCmdMock(t)
	defer gitCmd.MinimockFinish()

	tag := semantic.New(1, 2, 3)

	gitCmd.RunMock.When(
		[]string{"tag", "v1.2.3"},
		3*time.Second,
	).Then("", nil)

	gitCmd.RunMock.When(
		[]string{"push", "origin", "v1.2.3"},
		3*time.Second,
	).Then("", nil)

	tc := NewTagCreator(gitCmd)
	err := tc.CreateTag(tag)
	require.NoError(t, err)
}

func Test_TagCreator_errCreate(t *testing.T) {
	gitCmd := git.NewCmdMock(t)
	defer gitCmd.MinimockFinish()

	tag := semantic.New(1, 2, 3)

	gitCmd.RunMock.When(
		[]string{"tag", "v1.2.3"},
		3*time.Second,
	).Then("", errors.New("git create error"))

	tc := NewTagCreator(gitCmd)
	err := tc.CreateTag(tag)
	require.EqualError(t, err, "unable to create tag: git create error")
}

func Test_TagCreator_errPublish(t *testing.T) {
	gitCmd := git.NewCmdMock(t)
	defer gitCmd.MinimockFinish()

	tag := semantic.New(1, 2, 3)

	gitCmd.RunMock.When(
		[]string{"tag", "v1.2.3"},
		3*time.Second,
	).Then("", nil)

	gitCmd.RunMock.When(
		[]string{"push", "origin", "v1.2.3"},
		3*time.Second,
	).Then("", errors.New("git push error"))

	tc := NewTagCreator(gitCmd)
	err := tc.CreateTag(tag)
	require.EqualError(t, err, "unable to push tag: git push error")
}
