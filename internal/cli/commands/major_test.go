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

func Test_MajorCmd_commandInfo(t *testing.T) {
	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)

	majorCmd := NewMajorCmd(kit)

	name := majorCmd.Name()
	r.Equal(majorCmdName, name)

	synop := majorCmd.Synopsis()
	r.Equal(majorCmdSynopsis, synop)

	usage := majorCmd.Usage()
	r.Equal(majorCmdUsage, usage)
}

func Test_MajorCmd_Execute_normal(t *testing.T) {
	exp := "taggit: created tag v2.0.0\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	oldTag := semantic.New(1, 2, 3)
	newTag := semantic.New(2, 0, 0)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy{
			tags.NewTriple(1, 2, 3): {oldTag},
		}, nil,
	)
	mocks.tagCreator.CreateTagMock.When(newTag).Then(nil)
	mocks.tagPublisher.PublishMock.When(newTag).Then(nil)

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)
	majorCmd := NewMajorCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := majorCmd.Execute(ctx, fs)

	r.Equal(subcommands.ExitSuccess, status)
	r.Equal(exp, mocks.stdout.String())
	r.Equal("", mocks.stderr.String())
}

func Test_MajorCmd_Execute_noPrevious(t *testing.T) {
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
	majorCmd := NewMajorCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := majorCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}

func Test_MajorCmd_Execute_listErr(t *testing.T) {
	exp := "taggit: failure: some git error\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy(nil), errors.New("some git error"),
	)

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)
	majorCmd := NewMajorCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := majorCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}

func Test_MajorCmd_Execute_creatorErr(t *testing.T) {
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
		semantic.New(2, 0, 0),
	).Return(
		errors.New("some create error"),
	)

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)
	majorCmd := NewMajorCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := majorCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}

func Test_MajorCmd_Execute_publishErr(t *testing.T) {
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
		semantic.New(2, 0, 0),
	).Return(nil)

	mocks.tagPublisher.PublishMock.Expect(
		semantic.New(2, 0, 0),
	).Return(errors.New("some publish error"))

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)
	majorCmd := NewMajorCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := majorCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}
