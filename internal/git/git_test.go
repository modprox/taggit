package git

import (
	"errors"
	"testing"
	"time"

	"github.com/modprox/taggit/internal/git/gittest"
	"github.com/modprox/taggit/tags"
	"github.com/stretchr/testify/require"
)

func Test_ListTags_empty(t *testing.T) {
	cmd := &gittest.Cmd{}
	defer cmd.AssertExpectations(t)

	cmd.On(
		"Run",
		[]string{"tag", "-l"},
		10*time.Second,
	).Return("", nil).Once()

	repoTags, err := ListTags(cmd)
	require.NoError(t, err)

	require.Equal(t, []tags.Tag{}, repoTags)
}

func Test_ListTags_non_empty(t *testing.T) {
	cmd := &gittest.Cmd{}
	defer cmd.AssertExpectations(t)

	cmd.On(
		"Run",
		[]string{"tag", "-l"},
		10*time.Second,
	).Return("v0.2.0\nv1.1.2\n", nil).Once()

	repoTags, err := ListTags(cmd)
	require.NoError(t, err)

	require.Equal(t, []tags.Tag{
		tags.New(1, 1, 2),
		tags.New(0, 2, 0),
	}, repoTags)
}

func Test_ListTags_failed(t *testing.T) {
	cmd := &gittest.Cmd{}
	defer cmd.AssertExpectations(t)

	cmd.On(
		"Run",
		[]string{"tag", "-l"},
		10*time.Second,
	).Return(
		"",
		errors.New("git failed"),
	).Once()

	_, err := ListTags(cmd)
	require.Error(t, err)
}

func Test_CreateTag(t *testing.T) {
	cmd := &gittest.Cmd{}
	defer cmd.AssertExpectations(t)

	cmd.On(
		"Run",
		[]string{"tag", "v2.3.4"},
		3*time.Second,
	).Return("", nil).Once()

	cmd.On("Run",
		[]string{"push", "origin", "v2.3.4"},
		1*time.Minute,
	).Return("", nil).Once()

	err := CreateTag(cmd, tags.New(2, 3, 4))
	require.NoError(t, err)
}
