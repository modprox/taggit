package commands

import (
	"context"
	"errors"
	"flag"
	"testing"

	"github.com/google/subcommands"
	"github.com/stretchr/testify/require"
	"gophers.dev/pkgs/semantic"
	"oss.indeed.com/go/taggit/internal/tags"
)

func Test_PatchCmd_commandInfo(t *testing.T) {
	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)

	patchCmd := NewPatchCmd(kit)

	name := patchCmd.Name()
	r.Equal(patchCmdName, name)

	synop := patchCmd.Synopsis()
	r.Equal(patchCmdSynopsis, synop)

	usage := patchCmd.Usage()
	r.Equal(patchCmdUsage, usage)
}

func Test_PatchCmd_Execute_normal(t *testing.T) {
	exp := "taggit: created tag v1.2.4\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	oldTag := semantic.New(1, 2, 3)
	newTag := semantic.New(1, 2, 4)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy{
			tags.NewTriple(1, 2, 3): {oldTag},
		}, nil,
	)
	mocks.tagCreator.CreateTagMock.When(newTag).Then(nil)
	mocks.tagPublisher.PublishMock.When(newTag).Then(nil)

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)
	patchCmd := NewPatchCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := patchCmd.Execute(ctx, fs)

	r.Equal(subcommands.ExitSuccess, status)
	r.Equal(exp, mocks.stdout.String())
	r.Equal("", mocks.stderr.String())
}

func Test_PatchCmd_Execute_noPrevious(t *testing.T) {
	exp := `taggit: cannot increment tag because no previous tags exist
taggit: failure: no previous tags
`

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy(nil), nil, // no tags, no error
	)

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)
	patchCmd := NewPatchCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := patchCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}

func Test_PatchCmd_Execute_listErr(t *testing.T) {
	exp := "taggit: failure: some git error\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy(nil), errors.New("some git error"),
	)

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)
	patchCmd := NewPatchCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := patchCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}

func Test_PatchCmd_Execute_creatorErr(t *testing.T) {
	exp := "taggit: failure: some create error\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy{
			tags.NewTriple(1, 2, 3): {semantic.New(1, 2, 3)},
		}, nil,
	)

	mocks.tagCreator.CreateTagMock.Expect(
		semantic.New(1, 2, 4),
	).Return(
		errors.New("some create error"),
	)

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)
	patchCmd := NewPatchCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := patchCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}

func Test_PatchCmd_Execute_publishErr(t *testing.T) {
	exp := "taggit: failure: some publish error\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy{
			tags.NewTriple(1, 2, 3): {semantic.New(1, 2, 3)},
		}, nil,
	)

	mocks.tagCreator.CreateTagMock.Expect(
		semantic.New(1, 2, 4),
	).Return(nil)

	mocks.tagPublisher.PublishMock.Expect(
		semantic.New(1, 2, 4),
	).Return(errors.New("some publish error"))

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)
	patchCmd := NewPatchCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := patchCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}
