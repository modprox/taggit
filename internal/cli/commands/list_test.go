package commands

import (
	"context"
	"errors"
	"flag"
	"testing"

	"gophers.dev/pkgs/semantic"

	"github.com/google/subcommands"
	"github.com/stretchr/testify/require"

	"oss.indeed.com/go/taggit/internal/tags"
)

func Test_ListCmd_commandInfo(t *testing.T) {
	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)

	listCmd := NewListCmd(kit)

	name := listCmd.Name()
	r.Equal(listCmdName, name)

	synop := listCmd.Synopsis()
	r.Equal(listCmdSynopsis, synop)

	usage := listCmd.Usage()
	r.Equal(listCmdUsage, usage)
}

func Test_ListCmd_Execute_noTags(t *testing.T) {
	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy(nil), // no tags to parse
		nil,
	)

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)
	listCmd := NewListCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	status := listCmd.Execute(ctx, fs)

	r.Equal(subcommands.ExitSuccess, status)
}

func Test_ListCmd_Execute_someTags(t *testing.T) {
	exp := `v0.1.0 |= v0.1.0 v0.1.0-alpha1
v0.2.0 |= v0.2.0-rc1 v0.2.0-r1+linux v0.2.0-r1+darwin
`
	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		tags.Taxonomy{
			tags.NewTriple(0, 1, 0): []semantic.Tag{
				semantic.New(0, 1, 0),
				semantic.New2(0, 1, 0, "alpha1"),
			},
			tags.NewTriple(0, 2, 0): []semantic.Tag{
				semantic.New2(0, 2, 0, "rc1"),
				semantic.New3(0, 2, 0, "r1", "linux"),
				semantic.New3(0, 2, 0, "r1", "darwin"),
			},
		},
		nil,
	)

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)
	listCmd := NewListCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	status := listCmd.Execute(ctx, fs)

	r.Equal(subcommands.ExitSuccess, status)
	r.Equal(exp, mocks.stdout.String())
	r.Equal("", mocks.stderr.String())
}

func Test_ListCmd_Execute_listErr(t *testing.T) {
	exp := "taggit: failure: some git error\n"

	r := require.New(t)

	mocks := newMocks(t)
	defer mocks.assertions(t)

	mocks.tagLister.ListRepoTagsMock.Expect().Return(
		nil, errors.New("some git error"),
	)

	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)
	listCmd := NewListCmd(kit)

	ctx := context.Background()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	status := listCmd.Execute(ctx, fs)

	r.Equal(subcommands.ExitFailure, status)
	r.Equal("", mocks.stdout.String())
	r.Equal(exp, mocks.stderr.String())
}
