package commands

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/google/subcommands"

	"gophers.dev/pkgs/semantic"

	"oss.indeed.com/go/taggit/internal/cli"
	"oss.indeed.com/go/taggit/internal/cli/output"
	"oss.indeed.com/go/taggit/internal/tags"
)

const (
	listCmdName     = "list"
	listCmdSynopsis = "List tagged versions."
	listCmdUsage    = "list"
)

func NewListCmd(kit *Kit) subcommands.Command {
	return &listCmd{
		writer:    kit.writer,
		tagLister: kit.tagLister,
	}
}

type listCmd struct {
	writer    output.Writer
	tagLister cli.TagLister
}

func (lc *listCmd) Name() string {
	return listCmdName
}

func (lc *listCmd) Synopsis() string {
	return listCmdSynopsis
}

func (lc *listCmd) Usage() string {
	return listCmdUsage
}

func (lc *listCmd) SetFlags(fs *flag.FlagSet) {
	// no flags when listing versions, for now
}

func (lc *listCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if err := lc.execute(); err != nil {
		lc.writer.Errorf("failure: %v", err)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func (lc *listCmd) execute() error {
	tax, err := lc.tagLister.ListRepoTags()
	if err != nil {
		return err
	}

	lc.list(tax)
	return nil
}

func (lc *listCmd) list(groups tags.Taxonomy) {
	lc.writer.Tracef("listing tags in git repository")

	triples := groups.Bases()

	for _, triple := range triples {
		tagsOfTriple := groups[triple]
		line := outputLineForTriple(triple, tagsOfTriple)
		lc.writer.Directf("%s", line)
	}
}

func outputLineForTriple(triple tags.Triple, associated []semantic.Tag) string {
	asString := associatedList(associated)
	s := fmt.Sprintf("%s |= %s", triple, strings.Join(asString, " "))
	return s
}

func associatedList(associated []semantic.Tag) []string {
	var asStrings []string
	for _, aTag := range associated {
		asStrings = append(asStrings, aTag.String())
	}
	return asStrings
}
