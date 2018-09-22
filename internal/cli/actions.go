package cli

import (
	"bytes"
	"fmt"
	"io"

	"github.com/modprox/taggit/internal/git"
	"github.com/modprox/taggit/internal/tags"
)

// List writes each tag to w in order, separated by newline.
func List(w io.Writer, tags []tags.Tag) {
	var b bytes.Buffer
	for _, tag := range tags {
		b.WriteString(tag.String())
		b.WriteString("\n")
	}

	w.Write(b.Bytes())
}

// Zero creates the initial v0.0.0 tag.
func Zero(cmd git.Cmd, repoTags []tags.Tag) error {
	if len(repoTags) > 0 {
		return fmt.Errorf("refusing to generate zero tag (v0.0.0) when other semver tags already exist")
	}

	return git.CreateTag(cmd, tags.ZeroValue)
}

// Patch increments the patch level of the latest release by 1.
func Patch(cmd git.Cmd, repoTags []tags.Tag) error {
	if len(repoTags) < 1 {
		return fmt.Errorf("refusing to bump patch with no pre-existing tag")
	}

	latest := repoTags[0]

	return git.CreateTag(cmd, tags.Tag{
		Major: latest.Major,
		Minor: latest.Minor,
		Patch: latest.Patch + 1,
	})
}

// Minor increments the minor level of the latest release by 1
// and sets the patch level to 0.
func Minor(cmd git.Cmd, repoTags []tags.Tag) error {
	if len(repoTags) < 1 {
		return fmt.Errorf("refusing to bump minor with no pre-existing tag")
	}

	latest := repoTags[0]

	return git.CreateTag(cmd, tags.Tag{
		Major: latest.Major,
		Minor: latest.Minor + 1,
		Patch: 0,
	})
}

// Major increments the major level of the latest release by 1
// and sets the minor and patch levels to 0.
func Major(cmd git.Cmd, repoTags []tags.Tag) error {
	if len(repoTags) < 1 {
		return fmt.Errorf("refusing to bump major with no pre-existnig tag")
	}

	latest := repoTags[0]

	return git.CreateTag(cmd, tags.Tag{
		Major: latest.Major + 1,
		Minor: 0,
		Patch: 0,
	})
}
