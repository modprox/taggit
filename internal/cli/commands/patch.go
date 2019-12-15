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
	patchCmdName     = "patch"
	patchCmdSynopsis = "Create an incremented patch version"
	patchCmdUsage    = "patch [pre-release] -meta='build-metadata'"
)

func NewPatchCmd(kit *Kit) subcommands.Command {
	return &patchCmd{
		writer:       kit.writer,
		tagLister:    kit.tagLister,
		tagCreator:   kit.tagCreator,
		tagPublisher: kit.tagPublisher,
	}
}

type patchCmd struct {
	writer       output.Writer
	tagLister    cli.TagLister
	tagCreator   cli.TagCreator
	tagPublisher publish.Publisher
}

func (pc *patchCmd) Name() string {
	return patchCmdName
}

func (pc *patchCmd) Synopsis() string {
	return patchCmdSynopsis
}

func (pc *patchCmd) Usage() string {
	return patchCmdUsage
}

func (pc *patchCmd) SetFlags(fs *flag.FlagSet) {
	_ = fs.String("meta", "", "build metadata label")
}

func (pc *patchCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if err := pc.execute(tags.ExtractExtensions(fs)); err != nil {
		pc.writer.Errorf("failure: %v", err)
		return subcommands.ExitFailure
	}
	return subcommands.ExitSuccess
}

func (pc *patchCmd) execute(ext tags.Extensions) error {
	pc.writer.Tracef(
		"increment patch version, pre-release: %q, build-metadata: %q",
		ext.PreRelease, ext.BuildMetadta,
	)

	tax, err := pc.tagLister.ListRepoTags()
	if err != nil {
		return err
	}

	if exists := tags.HasPrevious(tax); !exists {
		pc.writer.Errorf("cannot increment tag because no previous tags exist")
		return ErrNoPreviousTags
	}

	latest := tax.Latest()
	next := tags.IncPatch(latest, ext)

	if err := pc.tagCreator.CreateTag(next); err != nil {
		return err
	}

	if err := pc.tagPublisher.Publish(next); err != nil {
		return err
	}

	pc.writer.Writef("created tag %s", next)
	return nil
}
