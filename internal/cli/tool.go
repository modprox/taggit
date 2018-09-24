package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/modprox/taggit/internal/git"
	"github.com/modprox/taggit/internal/publish"
	"github.com/modprox/taggit/internal/tags"
)

type Tool interface {
	List([]tags.Tag) error
	Zero([]tags.Tag) error
	Patch([]tags.Tag) error
	Minor([]tags.Tag) error
	Major([]tags.Tag) error
}

func NewTool(
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

func (t *tool) List(repoTags []tags.Tag) error {
	var b bytes.Buffer
	for _, tag := range repoTags {
		b.WriteString(tag.String())
		b.WriteString("\n")
	}
	_, err := t.output.Write(b.Bytes())
	return err
}

func (t *tool) Zero(repoTags []tags.Tag) error {
	if len(repoTags) > 0 {
		msg := "refusing to generate zero tag (v0.0.0) when other semver tags already exist"
		_ = t.write(msg)
		return errors.New("tags already exist")
	}

	if err := git.CreateTag(t.gitCmd, tags.Tag{
		Major: 0,
		Minor: 0,
		Patch: 0,
	}); err != nil {
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

func (t *tool) Patch(repoTags []tags.Tag) error {
	if len(repoTags) < 1 {
		msg := "refusing to bump patch with no pre-existing tag"
		t.write(msg)
		return errors.New("no tags exist yet")
	}

	latest := repoTags[0]
	newTag := tags.Tag{
		Major: latest.Major,
		Minor: latest.Minor,
		Patch: latest.Patch + 1,
	}

	if err := git.CreateTag(t.gitCmd, newTag); err != nil {
		msg := "failed to create tag: " + err.Error()
		t.write(msg)
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

func (t *tool) Minor(repoTags []tags.Tag) error {
	if len(repoTags) < 1 {
		msg := "refusing to bump minor with no pre-existing tag"
		_ = t.write(msg)
		return errors.New("no tags exist yet")
	}

	latest := repoTags[0]
	newTag := tags.Tag{
		Major: latest.Major,
		Minor: latest.Minor + 1,
		Patch: 0,
	}

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

func (t *tool) Major(repoTags []tags.Tag) error {
	if len(repoTags) < 1 {
		msg := "refusing to bump major with no pre-existing tag"
		_ = t.write(msg)
		return errors.New("no tags exist yet")
	}

	latest := repoTags[0]
	newTag := tags.Tag{
		Major: latest.Major + 1,
		Minor: 0,
		Patch: 0,
	}

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
