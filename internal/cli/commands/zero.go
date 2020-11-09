package commands

import (
	"context"
	"flag"

	"github.com/google/subcommands"
	"gophers.dev/pkgs/semantic"
	"oss.indeed.com/go/taggit/internal/cli"
	"oss.indeed.com/go/taggit/internal/cli/output"
	"oss.indeed.com/go/taggit/internal/publish"
	"oss.indeed.com/go/taggit/internal/tags"
)

const (
	zeroCmdName     = "zero"
	zeroCmdSynopsis = "Create initial v0.0.0 tag"
	zeroCmdUsage    = "zero"
)

func NewZeroCmd(kit *Kit) subcommands.Command {
	return &zeroCmd{
		writer:       kit.writer,
		tagLister:    kit.tagLister,
		tagCreator:   kit.tagCreator,
		tagPusher:    kit.tagPusher,
		tagPublisher: kit.tagPublisher,
	}
}

type zeroCmd struct {
	writer       output.Writer
	tagLister    cli.TagLister
	tagCreator   cli.TagCreator
	tagPusher    cli.TagPusher
	tagPublisher publish.Publisher
}

func (zc *zeroCmd) Name() string {
	return zeroCmdName
}

func (zc *zeroCmd) Synopsis() string {
	return zeroCmdSynopsis
}

func (zc *zeroCmd) Usage() string {
	return zeroCmdUsage
}

func (zc *zeroCmd) SetFlags(fs *flag.FlagSet) {
	// no build metadata when creating v0.0.0
}

func (zc *zeroCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if err := zc.execute(); err != nil {
		zc.writer.Errorf("failure: %v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func (zc *zeroCmd) execute() error {
	zc.writer.Tracef("create initial v0.0.0 tag")

	groups, err := zc.tagLister.ListRepoTags()
	if err != nil {
		return err
	}

	zero := semantic.New(0, 0, 0)

	if exists := tags.HasPrevious(groups); exists {
		zc.writer.Errorf("refusing to generate zero tag (%s) when other semver tags already exist", zero)
		return ErrPreviousTags
	}

	if err := zc.tagCreator.CreateTag(zero); err != nil {
		return err
	}

	if err := zc.tagPusher.PushTag(zero); err != nil {
		return err
	}

	if err := zc.tagPublisher.Publish(zero); err != nil {
		return err
	}

	zc.writer.Writef("created tag %s", zero)
	return nil
}
