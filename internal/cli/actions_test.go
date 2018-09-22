package cli

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/modprox/taggit/internal/git/gittest"
	"github.com/modprox/taggit/internal/tags"

	"github.com/stretchr/testify/require"
)

func Test_List(t *testing.T) {
	try := func(in []tags.Tag, exp string) {
		var w bytes.Buffer
		List(&w, in)
		output := w.String()
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
	cmd := &gittest.Cmd{}
	defer cmd.AssertExpectations(t)

	cmd.On(
		"Run",
		[]string{"tag", "v0.0.0"},
		3*time.Second).Return("", nil).Once()

	cmd.On(
		"Run",
		[]string{"push", "origin", "v0.0.0"},
		1*time.Minute).Return("", nil).Once()

	err := Zero(cmd, []tags.Tag{})
	require.NoError(t, err)
}

func Test_Zero_non_empty(t *testing.T) {
	cmd := &gittest.Cmd{}
	defer cmd.AssertExpectations(t)

	err := Zero(cmd, []tags.Tag{tags.New(0, 0, 1)})
	require.Error(t, err)
}

func Test_Zero_failure(t *testing.T) {
	cmd := &gittest.Cmd{}
	defer cmd.AssertExpectations(t)

	cmd.On(
		"Run",
		[]string{"tag", "v0.0.0"},
		3*time.Second).Return("", nil).Once()

	cmd.On(
		"Run",
		[]string{"push", "origin", "v0.0.0"},
		1*time.Minute).Return(
		"", errors.New("git is broken"),
	).Once()

	err := Zero(cmd, []tags.Tag{})
	require.Error(t, err)
}

func Test_Patch_empty(t *testing.T) {
	cmd := &gittest.Cmd{}
	defer cmd.AssertExpectations(t)

	err := Patch(cmd, []tags.Tag{})
	require.Error(t, err)
}

func Test_Patch_non_empty(t *testing.T) {
	cmd := &gittest.Cmd{}
	defer cmd.AssertExpectations(t)

	cmd.On(
		"Run",
		[]string{"tag", "v2.1.5"},
		3*time.Second,
	).Return("", nil).Once()

	cmd.On(
		"Run",
		[]string{"push", "origin", "v2.1.5"},
		1*time.Minute,
	).Return("", nil).Once()

	err := Patch(cmd, []tags.Tag{
		tags.New(2, 1, 4),
		tags.New(2, 0, 0),
		tags.New(1, 2, 3),
	})

	require.NoError(t, err)
}

func Test_Patch_failure(t *testing.T) {
	cmd := &gittest.Cmd{}
	defer cmd.AssertExpectations(t)

	cmd.On(
		"Run",
		[]string{"tag", "v2.1.5"},
		3*time.Second,
	).Return("", nil).Once()

	cmd.On(
		"Run",
		[]string{"push", "origin", "v2.1.5"},
		1*time.Minute,
	).Return(
		"", errors.New("git is broken"),
	).Once()

	err := Patch(cmd, []tags.Tag{
		tags.New(2, 1, 4),
		tags.New(2, 0, 0),
		tags.New(1, 2, 3),
	})

	require.Error(t, err)
}

func Test_Minor_empty(t *testing.T) {
	cmd := &gittest.Cmd{}
	defer cmd.AssertExpectations(t)

	err := Minor(cmd, []tags.Tag{})
	require.Error(t, err)
}

func Test_Minor_non_empty(t *testing.T) {
	cmd := &gittest.Cmd{}
	defer cmd.AssertExpectations(t)

	cmd.On(
		"Run",
		[]string{"tag", "v2.2.0"},
		3*time.Second,
	).Return("", nil).Once()

	cmd.On(
		"Run",
		[]string{"push", "origin", "v2.2.0"},
		1*time.Minute,
	).Return("", nil).Once()

	err := Minor(cmd, []tags.Tag{
		tags.New(2, 1, 4),
		tags.New(2, 0, 0),
		tags.New(1, 2, 3),
	})

	require.NoError(t, err)
}

func Test_Minor_failure(t *testing.T) {
	cmd := &gittest.Cmd{}
	defer cmd.AssertExpectations(t)

	cmd.On(
		"Run",
		[]string{"tag", "v2.2.0"},
		3*time.Second,
	).Return("", nil).Once()

	cmd.On(
		"Run",
		[]string{"push", "origin", "v2.2.0"},
		1*time.Minute,
	).Return(
		"", errors.New("git is broken"),
	).Once()

	err := Minor(cmd, []tags.Tag{
		tags.New(2, 1, 4),
		tags.New(2, 0, 0),
		tags.New(1, 2, 3),
	})

	require.Error(t, err)
}

func Test_Major_empty(t *testing.T) {
	cmd := &gittest.Cmd{}
	defer cmd.AssertExpectations(t)

	err := Major(cmd, []tags.Tag{})
	require.Error(t, err)
}

func Test_Major_non_empty(t *testing.T) {
	cmd := &gittest.Cmd{}
	defer cmd.AssertExpectations(t)

	cmd.On(
		"Run",
		[]string{"tag", "v3.0.0"},
		3*time.Second,
	).Return("", nil).Once()

	cmd.On(
		"Run",
		[]string{"push", "origin", "v3.0.0"},
		1*time.Minute,
	).Return("", nil).Once()

	err := Major(cmd, []tags.Tag{
		tags.New(2, 1, 4),
		tags.New(2, 0, 0),
		tags.New(1, 2, 3),
	})

	require.NoError(t, err)
}

func Test_Major_non_failure(t *testing.T) {
	cmd := &gittest.Cmd{}
	defer cmd.AssertExpectations(t)

	cmd.On(
		"Run",
		[]string{"tag", "v3.0.0"},
		3*time.Second,
	).Return("", nil).Once()

	cmd.On(
		"Run",
		[]string{"push", "origin", "v3.0.0"},
		1*time.Minute,
	).Return(
		"", errors.New("git is broken"),
	).Once()

	err := Major(cmd, []tags.Tag{
		tags.New(2, 1, 4),
		tags.New(2, 0, 0),
		tags.New(1, 2, 3),
	})

	require.Error(t, err)
}
