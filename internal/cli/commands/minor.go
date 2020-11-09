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
	minorCmdName     = "minor"
	minorCmdSynopsis = "Create an incremented minor version"
	minorCmdUsage    = "minor [pre-release] -meta='build-metadata'"
)

func NewMinorCmd(kit *Kit) subcommands.Command {
	return &minorCmd{
		writer:       kit.writer,
		tagLister:    kit.tagLister,
		tagCreator:   kit.tagCreator,
		tagPusher:    kit.tagPusher,
		tagPublisher: kit.tagPublisher,
	}
}

type minorCmd struct {
	writer       output.Writer
	tagLister    cli.TagLister
	tagCreator   cli.TagCreator
	tagPusher    cli.TagPusher
	tagPublisher publish.Publisher
}

func (mc *minorCmd) Name() string {
	return "minor"
}

func (mc *minorCmd) Synopsis() string {
	return "Create an incremented minor version"
}

func (mc *minorCmd) Usage() string {
	return "minor [pre-release] -meta='build-metadata'"
}

func (mc *minorCmd) SetFlags(fs *flag.FlagSet) {
	_ = fs.String("meta", "", "build metadata label")
}

func (mc *minorCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if err := mc.execute(tags.ExtractExtensions(fs)); err != nil {
		mc.writer.Errorf("failure: %v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func (mc *minorCmd) execute(ext tags.Extensions) error {
	mc.writer.Tracef(
		"increment minor version, pre-release: %q, build-metadata: %q",
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
	next := tags.IncMinor(latest, ext)

	if err := mc.tagCreator.CreateTag(next); err != nil {
		return err
	}

	if err := mc.tagPusher.PushTag(next); err != nil {
		return err
	}

	if err := mc.tagPublisher.Publish(next); err != nil {
		return err
	}

	mc.writer.Writef("created tag %s", next)
	return nil
}
