package cli

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/modprox/taggit/internal/git/gittest"
	"github.com/modprox/taggit/internal/publish"
	"github.com/modprox/taggit/internal/tags"

	"github.com/stretchr/testify/require"
)

func newTool(t *testing.T) (*tool, *bytes.Buffer, *gittest.Cmd) {
	var buf bytes.Buffer
	gitCmd := &gittest.Cmd{}
	publisher := publish.Discard()
	tool := NewTool(&buf, gitCmd, publisher).(*tool)
	return tool, &buf, gitCmd
}

func Test_List(t *testing.T) {

	try := func(in []tags.Tag, exp string) {
		tool, buf, gitCmd := newTool(t)
		defer gitCmd.AssertExpectations(t)

		err := tool.List(in)
		require.NoError(t, err)
		output := buf.String()
		require.Equal(t, exp, output)
	}

	try([]tags.Tag{}, "")

	try([]tags.Tag{
		tags.New(1, 2, 3),
	}, "v1.2.3\n")

	try([]tags.Tag{
		tags.New(0, 0, 1),
		tags.New(1, 0, 0),
		tags.New(0, 1, 0),
	}, "v0.0.1\nv1.0.0\nv0.1.0\n")
}

func Test_Zero_empty(t *testing.T) {
	tool, buf, gitCmd := newTool(t)
	defer gitCmd.AssertExpectations(t)

	gitCmd.On(
		"Run",
		[]string{"tag", "v0.0.0"},
		3*time.Second).Return("", nil).Once()

	gitCmd.On(
		"Run",
		[]string{"push", "origin", "v0.0.0"},
		1*time.Minute).Return("", nil).Once()

	err := tool.Zero([]tags.Tag{})
	require.NoError(t, err)
	require.Equal(t, "taggit: created tag: v0.0.0\ntaggit: published tag: v0.0.0\n", buf.String())
}

func Test_Zero_non_empty(t *testing.T) {
	tool, buf, gitCmd := newTool(t)
	defer gitCmd.AssertExpectations(t)

	err := tool.Zero([]tags.Tag{
		tags.New(0, 0, 1),
	})
	require.Error(t, err)
	require.Equal(t, "taggit: refusing to generate zero tag (v0.0.0) when other semver tags already exist\n", buf.String())
}

func Test_Zero_failure(t *testing.T) {
	tool, buf, gitCmd := newTool(t)
	defer gitCmd.AssertExpectations(t)

	gitCmd.On(
		"Run",
		[]string{"tag", "v0.0.0"},
		3*time.Second).Return("", nil).Once()

	gitCmd.On(
		"Run",
		[]string{"push", "origin", "v0.0.0"},
		1*time.Minute).Return(
		"", errors.New("git is broken"),
	).Once()

	err := tool.Zero([]tags.Tag{})
	require.Error(t, err)
	require.Equal(t, "taggit: failed to create tag: git is broken\n", buf.String())
}

func Test_Patch_empty(t *testing.T) {
	tool, buf, gitCmd := newTool(t)
	defer gitCmd.AssertExpectations(t)

	err := tool.Patch([]tags.Tag{})
	require.Error(t, err)
	require.Equal(t, "taggit: refusing to bump patch with no pre-existing tag\n", buf.String())
}

func Test_Patch_non_empty(t *testing.T) {
	tool, buf, gitCmd := newTool(t)
	defer gitCmd.AssertExpectations(t)

	gitCmd.On(
		"Run",
		[]string{"tag", "v2.1.5"},
		3*time.Second,
	).Return("", nil).Once()

	gitCmd.On(
		"Run",
		[]string{"push", "origin", "v2.1.5"},
		1*time.Minute,
	).Return("", nil).Once()

	err := tool.Patch([]tags.Tag{
		tags.New(2, 1, 4),
		tags.New(2, 0, 0),
		tags.New(1, 2, 3),
	})

	require.NoError(t, err)
	require.Equal(t, "taggit: created tag: v2.1.5\ntaggit: published tag: v2.1.5\n", buf.String())
}

func Test_Patch_failure(t *testing.T) {
	tool, buf, gitCmd := newTool(t)
	defer gitCmd.AssertExpectations(t)

	gitCmd.On(
		"Run",
		[]string{"tag", "v2.1.5"},
		3*time.Second,
	).Return("", nil).Once()

	gitCmd.On(
		"Run",
		[]string{"push", "origin", "v2.1.5"},
		1*time.Minute,
	).Return(
		"", errors.New("git is broken"),
	).Once()

	err := tool.Patch([]tags.Tag{
		tags.New(2, 1, 4),
		tags.New(2, 0, 0),
		tags.New(1, 2, 3),
	})

	require.Error(t, err)
	require.Equal(t, "taggit: failed to create tag: git is broken\n", buf.String())
}

func Test_Minor_empty(t *testing.T) {
	tool, buf, gitCmd := newTool(t)
	defer gitCmd.AssertExpectations(t)

	err := tool.Minor([]tags.Tag{})
	require.Error(t, err)
	require.Equal(t, "taggit: refusing to bump minor with no pre-existing tag\n", buf.String())
}

func Test_Minor_non_empty(t *testing.T) {
	tool, buf, gitCmd := newTool(t)
	defer gitCmd.AssertExpectations(t)

	gitCmd.On(
		"Run",
		[]string{"tag", "v2.2.0"},
		3*time.Second,
	).Return("", nil).Once()

	gitCmd.On(
		"Run",
		[]string{"push", "origin", "v2.2.0"},
		1*time.Minute,
	).Return("", nil).Once()

	err := tool.Minor([]tags.Tag{
		tags.New(2, 1, 4),
		tags.New(2, 0, 0),
		tags.New(1, 2, 3),
	})

	require.NoError(t, err)
	require.Equal(t, "taggit: created tag: v2.2.0\ntaggit: published tag: v2.2.0\n", buf.String())
}

func Test_Minor_failure(t *testing.T) {
	tool, buf, gitCmd := newTool(t)
	defer gitCmd.AssertExpectations(t)

	gitCmd.On(
		"Run",
		[]string{"tag", "v2.2.0"},
		3*time.Second,
	).Return("", nil).Once()

	gitCmd.On(
		"Run",
		[]string{"push", "origin", "v2.2.0"},
		1*time.Minute,
	).Return(
		"", errors.New("git is broken"),
	).Once()

	err := tool.Minor([]tags.Tag{
		tags.New(2, 1, 4),
		tags.New(2, 0, 0),
		tags.New(1, 2, 3),
	})

	require.Error(t, err)
	require.Equal(t, "taggit: failed to create tag: git is broken\n", buf.String())
}

func Test_Major_empty(t *testing.T) {
	tool, buf, gitCmd := newTool(t)
	defer gitCmd.AssertExpectations(t)

	err := tool.Major([]tags.Tag{})
	require.Error(t, err)
	require.Equal(t, "taggit: refusing to bump major with no pre-existing tag\n", buf.String())
}

func Test_Major_non_empty(t *testing.T) {
	tool, buf, gitCmd := newTool(t)
	defer gitCmd.AssertExpectations(t)

	gitCmd.On(
		"Run",
		[]string{"tag", "v3.0.0"},
		3*time.Second,
	).Return("", nil).Once()

	gitCmd.On(
		"Run",
		[]string{"push", "origin", "v3.0.0"},
		1*time.Minute,
	).Return("", nil).Once()

	err := tool.Major([]tags.Tag{
		tags.New(2, 1, 4),
		tags.New(2, 0, 0),
		tags.New(1, 2, 3),
	})

	require.NoError(t, err)
	require.Equal(t, "taggit: created tag: v3.0.0\ntaggit: published tag: v3.0.0\n", buf.String())
}

func Test_Major_failure(t *testing.T) {
	tool, buf, gitCmd := newTool(t)
	defer gitCmd.AssertExpectations(t)

	gitCmd.On(
		"Run",
		[]string{"tag", "v3.0.0"},
		3*time.Second,
	).Return("", nil).Once()

	gitCmd.On(
		"Run",
		[]string{"push", "origin", "v3.0.0"},
		1*time.Minute,
	).Return(
		"", errors.New("git is broken"),
	).Once()

	err := tool.Major([]tags.Tag{
		tags.New(2, 1, 4),
		tags.New(2, 0, 0),
		tags.New(1, 2, 3),
	})

	require.Error(t, err)
	require.Equal(t, "taggit: failed to create tag: git is broken\n", buf.String())
}
