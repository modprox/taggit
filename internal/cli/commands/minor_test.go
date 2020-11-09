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

func Test_MinorCmd_commandInfo(t *testing.T) {
	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	majorCmd := NewMinorCmd(newKit(mocks))

	name := majorCmd.Name()
	r.Equal(minorCmdName, name)

	synop := majorCmd.Synopsis()
	r.Equal(minorCmdSynopsis, synop)

	usage := majorCmd.Usage()
	r.Equal(minorCmdUsage, usage)
}

func Test_MinorCmd_Execute_normal(t *testing.T) {
	exp := "taggit: created tag v1.3.0\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	oldTag := semantic.New(1, 2, 3)
	newTag := semantic.New(1, 3, 0)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy{
			tags.NewTriple(1, 2, 3): {oldTag},
		}, nil,
	)
	mocks.tagCreator.CreateTagMock.When(newTag).Then(nil)
	mocks.tagPusher.PushTagMock.When(newTag).Then(nil)
	mocks.tagPublisher.PublishMock.When(newTag).Then(nil)

	minorCmd := NewMinorCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := minorCmd.Execute(ctx, fs)

	r.Equal(subcommands.ExitSuccess, status)
	r.Equal(exp, mocks.stdout.String())
	r.Equal("", mocks.stderr.String())
}

func Test_MinorCmd_Execute_noPrevious(t *testing.T) {
	exp := `taggit: cannot increment tag because no previous tags exist
taggit: failure: no previous tags
`

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy(nil), nil, // no tags, no error
	)

	minorCmd := NewMinorCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := minorCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}

func Test_MinorCmd_Execute_listErr(t *testing.T) {
	exp := "taggit: failure: some git error\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy(nil), errors.New("some git error"),
	)

	minorCmd := NewMinorCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := minorCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}

func Test_MinorCmd_Execute_creatorErr(t *testing.T) {
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
		semantic.New(1, 3, 0),
	).Return(
		errors.New("some create error"),
	)

	minorCmd := NewMinorCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := minorCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}

func Test_MinorCmd_Execute_publishErr(t *testing.T) {
	exp := "taggit: failure: some publish error\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy{
			tags.NewTriple(1, 2, 3): {semantic.New(1, 2, 3)},
		}, nil,
	)
	newTag := semantic.New(1, 3, 0)
	mocks.tagCreator.CreateTagMock.Expect(newTag).Return(nil)
	mocks.tagPusher.PushTagMock.Expect(newTag).Return(nil)
	mocks.tagPublisher.PublishMock.Expect(newTag).Return(errors.New("some publish error"))

	minorCmd := NewMinorCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := minorCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}
