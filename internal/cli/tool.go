package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/modprox/taggit/internal/updates"

	"github.com/modprox/taggit/internal/git"
	"github.com/modprox/taggit/internal/publish"
	"github.com/modprox/taggit/tags"
)

type Tool interface {
	List(Groups) error
	Updates(string) error
	Zero(Groups) error
	Patch(Groups, string) error
	Minor(Groups, string) error
	Major(Groups, string) error
}

func newTool(
	output io.Writer,
	gitCmd git.Cmd,
	publisher publish.Publisher,
) Tool {
	return &tool{
		output:    output,
		gitCmd:    gitCmd,
		publisher: publisher,
	}
}

type tool struct {
	output    io.Writer
	gitCmd    git.Cmd
	publisher publish.Publisher
}

func (t *tool) write(s string) error {
	_, err := t.output.Write([]byte("taggit: " + s + "\n"))
	return err
}

func (t *tool) List(repoTags Groups) error {
	var b bytes.Buffer

	// iterate the tags in semver order
	for _, base := range repoTags.Bases() {
		tagsOfBase := repoTags[base]

		b.WriteString(base.String())
		if tagsOfBase[0].IsBase() {
			b.WriteString(" ✓ |")
		} else {
			b.WriteString("   |")
		}

		for _, tag := range tagsOfBase {
			b.WriteString(" ")
			if tag.IsBase() {
				b.WriteString(tag.String())
			} else {
				b.WriteString(tag.String())
			}
		}

		b.WriteString("\n")
	}

	_, err := t.output.Write(b.Bytes())
	return err
}

func (t *tool) Updates(extension string) error {
	lister := updates.NewLister(updates.Options{
		Timeout: 10 * time.Second,
	})

	available, err := lister.List()
	if err != nil {
		return err
	}

	for _, up := range available {
		_ = t.write(up.String() + "\n")
	}

	return nil
}

func (t *tool) Zero(repoTags Groups) error {
	if len(repoTags) > 0 {
		msg := "refusing to generate zero tag (v0.0.0) when other semver tags already exist"
		_ = t.write(msg)
		return errors.New("tags already exist")
	}

	if err := git.CreateTag(t.gitCmd, tags.New(
		0, 0, 0,
	)); err != nil {
		msg := "failed to create tag: " + err.Error()
		_ = t.write(msg)
		return err
	}

	if err := t.write("created tag: v0.0.0"); err != nil {
		return err
	}

	if err := t.publisher.Publish("v0.0.0"); err != nil {
		return err
	}

	return t.write("published tag: v0.0.0")
}

func (t *tool) Patch(repoTags Groups, ext string) error {
	if len(repoTags) < 1 {
		msg := "refusing to bump patch with no pre-existing tag"
		_ = t.write(msg)
		return errors.New("no tags exist yet")
	}

	latest := repoTags.Latest()

	// compute the patch level - if latest is a base version,
	// we must bump the patch level. if latest is a pre-release
	// we continue forward on the current patch level
	patch := latest.Patch
	if latest.IsBase() {
		patch++
	}

	newTag := tags.New2(
		latest.Major,
		latest.Minor,
		patch,
		ext,
	)

	if err := git.CreateTag(t.gitCmd, newTag); err != nil {
		msg := "failed to create tag: " + err.Error()
		_ = t.write(msg)
		return err
	}

	msg := fmt.Sprintf("created tag: %s", newTag)
	if err := t.write(msg); err != nil {
		return err
	}

	if err := t.publisher.Publish(newTag.String()); err != nil {
		return err
	}

	return t.write("published tag: " + newTag.String())
}

func (t *tool) Minor(repoTags Groups, ext string) error {
	if len(repoTags) < 1 {
		msg := "refusing to bump minor with no pre-existing tag"
		_ = t.write(msg)
		return errors.New("no tags exist yet")
	}

	latest := repoTags.Latest()

	// compute the minor & patch level - if latest is a base version,
	// we bump the minor level and the patch goes to zero. if latest
	// is a pre-release version, we continue on the current minor and
	// patch level.

	minor := latest.Minor
	patch := latest.Patch
	if latest.IsBase() {
		minor++
		patch = 0
	}

	newTag := tags.New2(
		latest.Major,
		minor,
		patch,
		ext,
	)

	if err := git.CreateTag(t.gitCmd, newTag); err != nil {
		msg := "failed to create tag: " + err.Error()
		_ = t.write(msg)
		return err
	}

	msg := fmt.Sprintf("created tag: %s", newTag)
	if err := t.write(msg); err != nil {
		return err
	}

	if err := t.publisher.Publish(newTag.String()); err != nil {
		return err
	}

	return t.write("published tag: " + newTag.String())
}

func (t *tool) Major(repoTags Groups, ext string) error {
	if len(repoTags) < 1 {
		msg := "refusing to bump major with no pre-existing tag"
		_ = t.write(msg)
		return errors.New("no tags exist yet")
	}

	latest := repoTags.Latest()

	// compute the major, minor, and patch level - if latest is a base
	// version, we bump the major version, and the minor and patch levels
	// go to zero. if latest is a pre-release we continue forward on
	// the current levels

	major := latest.Major
	minor := latest.Minor
	patch := latest.Patch
	if latest.IsBase() {
		major++
		minor = 0
		patch = 0
	}

	newTag := tags.New2(
		major,
		minor,
		patch,
		ext,
	)

	if err := git.CreateTag(t.gitCmd, newTag); err != nil {
		msg := "failed to create tag: " + err.Error()
		_ = t.write(msg)
		return err
	}

	msg := fmt.Sprintf("created tag: %s", newTag)
	if err := t.write(msg); err != nil {
		return err
	}

	if err := t.publisher.Publish(newTag.String()); err != nil {
		return err
	}

	return t.write("published tag: " + newTag.String())
}
