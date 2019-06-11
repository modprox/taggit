package git

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"oss.indeed.com/go/taggit/tags"
)

const (
	timeout3s  = 3 * time.Second
	timeout10S = 10 * time.Second
	timeout1M  = 1 * time.Minute
)

func Test_ListTags_empty(t *testing.T) {
	cmd := NewCmdMock(t)
	defer cmd.MinimockFinish()

	cmd.RunMock.When(
		[]string{"tag", "-l"}, 10*time.Second,
	).Then("", nil)

	tagsContent, err := ListTags(cmd)
	require.NoError(t, err)

	require.Equal(t, "", tagsContent)
}

func Test_ListTags_non_empty(t *testing.T) {
	cmd := NewCmdMock(t)
	defer cmd.MinimockFinish()

	cmd.RunMock.When(
		[]string{"tag", "-l"}, timeout10S,
	).Then("v0.2.0\nv1.1.2\n", nil)

	tagsContent, err := ListTags(cmd)
	require.NoError(t, err)

	require.Equal(t, "v0.2.0\nv1.1.2\n", tagsContent)
}

func Test_ListTags_failed(t *testing.T) {
	cmd := NewCmdMock(t)
	defer cmd.MinimockFinish()

	cmd.RunMock.When(
		[]string{"tag", "-l"}, timeout10S,
	).Then("", errors.New("git failed"))

	_, err := ListTags(cmd)
	require.Error(t, err)
}

func Test_CreateTag(t *testing.T) {
	cmd := NewCmdMock(t)
	defer cmd.MinimockFinish()

	cmd.RunMock.When(
		[]string{"tag", "v2.3.4"}, timeout3s,
	).Then("", nil)

	cmd.RunMock.When(
		[]string{"push", "origin", "v2.3.4"}, timeout1M,
	).Then("", nil)

	err := CreateTag(cmd, tags.New(2, 3, 4))
	require.NoError(t, err)
}
