package commands

import (
	"context"
	"flag"

	"github.com/google/subcommands"

	"oss.indeed.com/go/taggit/internal/cli"
	"oss.indeed.com/go/taggit/internal/cli/output"
	"oss.indeed.com/go/taggit/internal/publish"
	"oss.indeed.com/go/taggit/internal/tags"
)

const (
	majorCmdName     = "major"
	majorCmdSynopsis = "Create an incremented major version"
	majorCmdUsage    = "major [pre-release] -meta='build-metadata'"
)

func NewMajorCmd(kit *Kit) subcommands.Command {
	return &majorCmd{
		writer:       kit.writer,
		tagLister:    kit.tagLister,
		tagCreator:   kit.tagCreator,
		tagPublisher: kit.tagPublisher,
	}
}

type majorCmd struct {
	writer       output.Writer
	tagLister    cli.TagLister
	tagCreator   cli.TagCreator
	tagPublisher publish.Publisher
}

func (mc *majorCmd) Name() string {
	return majorCmdName
}

func (mc *majorCmd) Synopsis() string {
	return majorCmdSynopsis
}

func (mc *majorCmd) Usage() string {
	return majorCmdUsage
}

func (mc *majorCmd) SetFlags(fs *flag.FlagSet) {
	_ = fs.String("meta", "", "build metadata label")
}

func (mc *majorCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if err := mc.execute(tags.ExtractExtensions(fs)); err != nil {
		mc.writer.Errorf("failure: %v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func (mc *majorCmd) execute(ext tags.Extensions) error {
	mc.writer.Tracef(
		"increment major version, pre-release: %q, build-metadata: %q",
		ext.PreRelease, ext.BuildMetadta,
	)

	groups, err := mc.tagLister.ListRepoTags()
	if err != nil {
		return err
	}

	if exists := tags.HasPrevious(groups); !exists {
		mc.writer.Errorf("cannot increment tag because no previous tags exist")
		return ErrNoPreviousTags
	}

	latest := groups.Latest()
	next := tags.IncMajor(latest)
	next.PreRelease = ext.PreRelease
	next.BuildMetadata = ext.BuildMetadta

	if err := mc.tagCreator.CreateTag(next); err != nil {
		return err
	}

	if err := mc.tagPublisher.Publish(next); err != nil {
		return err
	}

	mc.writer.Writef("created tag %s", next)
	return nil
}
