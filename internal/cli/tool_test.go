package cli

import (
	"bytes"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/require"

	"oss.indeed.com/go/taggit/internal/git"
	"oss.indeed.com/go/taggit/internal/publish"
)

func newTestTool(t *testing.T) (*tool, *bytes.Buffer, *git.CmdMock) {
	var buf bytes.Buffer
	gitCmd := git.NewCmdMock(t)
	publisher := publish.Discard()
	tool := newTool(&buf, gitCmd, publisher).(*tool)
	return tool, &buf, gitCmd
}

const (
	tagsEmpty    = ""
	tagsOne      = "v1.2.3\n"
	tagsOneAlpha = "v1.2.3-alpha1"
	tagsThree    = "v0.0.1\nv1.0.0\nv0.1.0\n"
	tagsThreeMix = "v1.1.1-beta3\nv1.1.1-beta4\nv1.1.2\n"
)

func Test_List(t *testing.T) {

	try := func(in string, exp string) {
		tool, buf, gitCmd := newTestTool(t)
		defer gitCmd.MinimockFinish()

		groups, err := Parse(in)
		require.NoError(t, err)

		err = tool.List(groups)
		require.NoError(t, err)
		output := buf.String()
		require.Equal(t, exp, output)
	}

	// empty
	try(tagsEmpty, "")

	// one basic
	try(tagsOne, "v1.2.3 ✓ | v1.2.3\n")

	// one extended
	try(tagsOneAlpha, "v1.2.3   | v1.2.3-alpha1\n")

	// three basic
	try(tagsThree, "v0.0.1 ✓ | v0.0.1\nv0.1.0 ✓ | v0.1.0\nv1.0.0 ✓ | v1.0.0\n")

	// three mixed
	try(tagsThreeMix, "v1.1.1   | v1.1.1-beta4 v1.1.1-beta3\nv1.1.2 ✓ | v1.1.2\n")
}

func tryParse(t *testing.T, in string) Groups {
	groups, err := Parse(in)
	require.NoError(t, err)
	return groups
}

const (
	timeout3S = 3 * time.Second
	timeout1M = 1 * time.Minute
)

func Test_Zero_empty(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	gitCmd.RunMock.When(
		[]string{"tag", "v0.0.0"}, timeout3S,
	).Then(
		"", nil,
	)

	gitCmd.RunMock.When(
		[]string{"push", "origin", "v0.0.0"}, timeout1M,
	).Then(
		"", nil,
	)

	groups := tryParse(t, tagsEmpty)
	err := tool.Zero(groups)
	require.NoError(t, err)
	require.Equal(t, "taggit: created tag: v0.0.0\ntaggit: published tag: v0.0.0\n", buf.String())
}

func Test_Zero_non_empty(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	groups := tryParse(t, tagsOne)
	err := tool.Zero(groups)
	require.Error(t, err)
	require.Equal(t, "taggit: refusing to generate zero tag (v0.0.0) when other semver tags already exist\n", buf.String())
}

func Test_Zero_failure(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	gitCmd.RunMock.When(
		[]string{"tag", "v0.0.0"}, timeout3S,
	).Then(
		"", nil,
	)

	gitCmd.RunMock.When(
		[]string{"push", "origin", "v0.0.0"}, timeout1M,
	).Then(
		"", errors.New("git is broken"),
	)

	groups := tryParse(t, tagsEmpty)

	err := tool.Zero(groups)
	require.Error(t, err)
	require.Equal(t, "taggit: failed to create tag: git is broken\n", buf.String())
}

func Test_Patch_empty(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	groups := tryParse(t, tagsEmpty)

	err := tool.Patch(groups, "")
	require.Error(t, err)
	require.Equal(t, "taggit: refusing to bump patch with no pre-existing tag\n", buf.String())
}

func Test_Patch_non_empty(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	gitCmd.RunMock.When(
		[]string{"tag", "v2.1.5"}, timeout3S,
	).Then("", nil)

	gitCmd.RunMock.When(
		[]string{"push", "origin", "v2.1.5"}, timeout1M,
	).Then("", nil)

	groups := tryParse(t, "v2.1.4\nv2.0.0\nv1.2.3\n")

	err := tool.Patch(groups, "")

	require.NoError(t, err)
	require.Equal(t, "taggit: created tag: v2.1.5\ntaggit: published tag: v2.1.5\n", buf.String())
}

func Test_Patch_extended(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	gitCmd.RunMock.When(
		[]string{"tag", "v2.1.5-alpha1"}, timeout3S,
	).Then("", nil)

	gitCmd.RunMock.When(
		[]string{"push", "origin", "v2.1.5-alpha1"}, timeout1M,
	).Then("", nil)

	groups := tryParse(t, "v2.1.4\nv2.0.0\nv1.2.3\n")

	err := tool.Patch(groups, "alpha1") // w/ extension

	require.NoError(t, err)
	require.Equal(t, "taggit: created tag: v2.1.5-alpha1\ntaggit: published tag: v2.1.5-alpha1\n", buf.String())
}

func Test_Patch_promote(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	gitCmd.RunMock.When(
		[]string{"tag", "v2.1.5"}, timeout3S,
	).Then("", nil)

	gitCmd.RunMock.When(
		[]string{"push", "origin", "v2.1.5"}, timeout1M,
	).Then("", nil)

	// v2.1.5-alpha3 incurs promotion of 2.1.5 base
	groups := tryParse(t, "v2.1.5-alpha3\nv2.1.5-alpha2\nv2.1.4\n")

	err := tool.Patch(groups, "")

	require.NoError(t, err)
	require.Equal(t, "taggit: created tag: v2.1.5\ntaggit: published tag: v2.1.5\n", buf.String())
}

func Test_Patch_failure(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	gitCmd.RunMock.When(
		[]string{"tag", "v2.1.5"},
		timeout3S,
	).Then("", nil)

	gitCmd.RunMock.When(
		[]string{"push", "origin", "v2.1.5"}, timeout1M,
	).Then("", errors.New("git is broken"))

	groups := tryParse(t, "v2.1.4\nv2.0.0\nv1.2.3\n")

	err := tool.Patch(groups, "")

	require.Error(t, err)
	require.Equal(t, "taggit: failed to create tag: git is broken\n", buf.String())
}

func Test_Minor_empty(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	groups := tryParse(t, tagsEmpty)

	err := tool.Minor(groups, "")
	require.Error(t, err)
	require.Equal(t, "taggit: refusing to bump minor with no pre-existing tag\n", buf.String())
}

func Test_Minor_non_empty(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	gitCmd.RunMock.When(
		[]string{"tag", "v2.2.0"}, timeout3S,
	).Then("", nil)

	gitCmd.RunMock.When(
		[]string{"push", "origin", "v2.2.0"}, timeout1M,
	).Then("", nil)

	groups := tryParse(t, "v2.1.4\nv2.0.0\nv1.2.3\n")

	err := tool.Minor(groups, "")

	require.NoError(t, err)
	require.Equal(t, "taggit: created tag: v2.2.0\ntaggit: published tag: v2.2.0\n", buf.String())
}

func Test_Minor_extended(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	gitCmd.RunMock.When(
		[]string{"tag", "v2.2.0-beta2"}, timeout3S,
	).Then("", nil)

	gitCmd.RunMock.When(
		[]string{"push", "origin", "v2.2.0-beta2"}, timeout1M,
	).Then("", nil)

	groups := tryParse(t, "v2.1.4\nv2.0.0\nv1.2.3\n")

	err := tool.Minor(groups, "beta2")

	require.NoError(t, err)
	require.Equal(t, "taggit: created tag: v2.2.0-beta2\ntaggit: published tag: v2.2.0-beta2\n", buf.String())
}

func Test_Minor_promote(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	gitCmd.RunMock.When(
		[]string{"tag", "v2.2.0"}, timeout3S,
	).Then("", nil)

	gitCmd.RunMock.When(
		[]string{"push", "origin", "v2.2.0"}, timeout1M,
	).Then("", nil)

	// v2.2.0-beta incurs promotion of v2.2.0 base
	groups := tryParse(t, "v2.2.0-beta\nv2.2.0-alpha2\nv2.2.0-alpha1\n")

	err := tool.Minor(groups, "")

	require.NoError(t, err)
	require.Equal(t, "taggit: created tag: v2.2.0\ntaggit: published tag: v2.2.0\n", buf.String())
}

func Test_Minor_failure(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	gitCmd.RunMock.When(
		[]string{"tag", "v2.2.0"}, timeout3S,
	).Then("", nil)

	gitCmd.RunMock.When(
		[]string{"push", "origin", "v2.2.0"}, timeout1M,
	).Then("", errors.New("git is broken"))

	groups := tryParse(t, "v2.1.4\nv2.0.0\nv1.2.3\n")

	err := tool.Minor(groups, "")

	require.Error(t, err)
	require.Equal(t, "taggit: failed to create tag: git is broken\n", buf.String())
}

func Test_Major_empty(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	groups := tryParse(t, tagsEmpty)

	err := tool.Major(groups, "")
	require.Error(t, err)
	require.Equal(t, "taggit: refusing to bump major with no pre-existing tag\n", buf.String())
}

func Test_Major_non_empty(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	gitCmd.RunMock.When(
		[]string{"tag", "v3.0.0"}, timeout3S,
	).Then("", nil)

	gitCmd.RunMock.When(
		[]string{"push", "origin", "v3.0.0"}, timeout1M,
	).Then("", nil)

	groups := tryParse(t, "v2.1.4\nv2.0.0\nv1.2.3\n")

	err := tool.Major(groups, "")

	require.NoError(t, err)
	require.Equal(t, "taggit: created tag: v3.0.0\ntaggit: published tag: v3.0.0\n", buf.String())
}

func Test_Major_extended(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	gitCmd.RunMock.When(
		[]string{"tag", "v3.0.0-rc1"}, timeout3S,
	).Then("", nil)

	gitCmd.RunMock.When(
		[]string{"push", "origin", "v3.0.0-rc1"}, timeout1M,
	).Then("", nil)

	groups := tryParse(t, "v2.1.4\nv2.0.0\nv1.2.3\n")

	err := tool.Major(groups, "rc1")

	require.NoError(t, err)
	require.Equal(t, "taggit: created tag: v3.0.0-rc1\ntaggit: published tag: v3.0.0-rc1\n", buf.String())
}

func Test_Major_promotion(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	gitCmd.RunMock.When(
		[]string{"tag", "v3.0.0"}, timeout3S,
	).Then("", nil)

	gitCmd.RunMock.When(
		[]string{"push", "origin", "v3.0.0"}, timeout1M,
	).Then("", nil)

	groups := tryParse(t, "v3.0.0-rc2\nv3.0.0-rc1\nv3.0.0-beta3\nv2.0.0\nv1.2.3\n")

	err := tool.Major(groups, "")

	require.NoError(t, err)
	require.Equal(t, "taggit: created tag: v3.0.0\ntaggit: published tag: v3.0.0\n", buf.String())
}

func Test_Major_failure(t *testing.T) {
	tool, buf, gitCmd := newTestTool(t)
	defer gitCmd.MinimockFinish()

	gitCmd.RunMock.When(
		[]string{"tag", "v3.0.0"}, timeout3S,
	).Then("", nil)

	gitCmd.RunMock.When(
		[]string{"push", "origin", "v3.0.0"}, timeout1M,
	).Then("", errors.New("git is broken"))

	groups := tryParse(t, "v2.1.4\nv2.0.0\nv1.2.3\n")

	err := tool.Major(groups, "")

	require.Error(t, err)
	require.Equal(t, "taggit: failed to create tag: git is broken\n", buf.String())
}
