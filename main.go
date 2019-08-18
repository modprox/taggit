// Command taggit publishes new versions of Go modules.
package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"

	"oss.indeed.com/go/taggit/internal/cli"
	"oss.indeed.com/go/taggit/internal/cli/commands"
	"oss.indeed.com/go/taggit/internal/cli/output"
	"oss.indeed.com/go/taggit/internal/git"
	"oss.indeed.com/go/taggit/internal/publish"
)

func main() {
	writer := output.NewWriter(os.Stdout, os.Stderr)
	gitCmd := git.New("git")
	tagLister := cli.NewTagLister(gitCmd)
	tagCreator := cli.NewTagCreator(gitCmd)
	tagPublisher := publish.FromEnv(publish.RegistryEnv, writer)
	kit := commands.NewKit(writer, tagLister, tagCreator, tagPublisher)

	listCmd := commands.NewListCmd(kit)
	zeroCmd := commands.NewZeroCmd(kit)
	patchCmd := commands.NewPatchCmd(kit)
	minorCmd := commands.NewMinorCmd(kit)
	majorCmd := commands.NewMajorCmd(kit)

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	subs := subcommands.NewCommander(fs, "")
	subs.Register(subs.HelpCommand(), "")
	subs.Register(subs.FlagsCommand(), "")
	subs.Register(listCmd, "")
	subs.Register(zeroCmd, "")
	subs.Register(patchCmd, "")
	subs.Register(minorCmd, "")
	subs.Register(majorCmd, "")

	if err := fs.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	ctx := context.Background()
	rc := subs.Execute(ctx, fs.Args())
	os.Exit(int(rc))
}
