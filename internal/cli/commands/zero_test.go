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

func Test_ZeroCmd_commandInfo(t *testing.T) {
	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	majorCmd := NewZeroCmd(newKit(mocks))

	name := majorCmd.Name()
	r.Equal(zeroCmdName, name)

	synop := majorCmd.Synopsis()
	r.Equal(zeroCmdSynopsis, synop)

	usage := majorCmd.Usage()
	r.Equal(zeroCmdUsage, usage)
}

func Test_ZeroCmd_Execute_normal(t *testing.T) {
	exp := "taggit: created tag v0.0.0\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	zeroTag := semantic.New(0, 0, 0)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		nil, nil, // no tags
	)
	mocks.tagCreator.CreateTagMock.When(zeroTag).Then(nil)
	mocks.tagPusher.PushTagMock.When(zeroTag).Then(nil)
	mocks.tagPublisher.PublishMock.When(zeroTag).Then(nil)

	zeroCmd := NewZeroCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := zeroCmd.Execute(ctx, fs)

	r.Equal(subcommands.ExitSuccess, status)
	r.Equal(exp, mocks.stdout.String())
	r.Equal("", mocks.stderr.String())
}

func Test_ZeroCmd_Execute_hasPrevious(t *testing.T) {
	exp := "refusing to generate zero tag (v0.0.0) when other semver tags already exist\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	oldTag := semantic.New(1, 2, 3)
	// zeroTag := semantic.New(0, 0, 0)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy{
			tags.NewTriple(1, 2, 3): {oldTag},
		}, nil,
	)
	// mocks.tagCreator.CreateTagMock.When(zeroTag).Then(nil)
	// mocks.tagPublisher.PublishMock.When(zeroTag).Then(nil)

	zeroCmd := NewZeroCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := zeroCmd.Execute(ctx, fs)

	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Contains(mocks.stderr.String(), exp)
}

func Test_ZeroCmd_Execute_listErr(t *testing.T) {
	exp := "taggit: failure: some git error\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy(nil), errors.New("some git error"),
	)

	zeroCmd := NewZeroCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := zeroCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}

func Test_ZeroCmd_Execute_creatorErr(t *testing.T) {
	exp := "taggit: failure: some create error\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		nil, nil,
	)

	mocks.tagCreator.CreateTagMock.Expect(
		semantic.New(0, 0, 0),
	).Return(
		errors.New("some create error"),
	)

	zeroCmd := NewZeroCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := zeroCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}

func Test_ZeroCmd_Execute_publishErr(t *testing.T) {
	exp := "taggit: failure: some publish error\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		nil, nil,
	)
	newTag := semantic.New(0, 0, 0)
	mocks.tagCreator.CreateTagMock.Expect(newTag).Return(nil)
	mocks.tagPusher.PushTagMock.Expect(newTag).Return(nil)
	mocks.tagPublisher.PublishMock.Expect(newTag).Return(errors.New("some publish error"))

	zeroCmd := NewZeroCmd(newKit(mocks))

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	_ = fs.String("meta", "", "usage")

	status := zeroCmd.Execute(ctx, fs)
	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}
