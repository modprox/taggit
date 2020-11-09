package commands

import (
	"oss.indeed.com/go/taggit/internal/cli"
	"oss.indeed.com/go/taggit/internal/cli/output"
	"oss.indeed.com/go/taggit/internal/publish"
)

func NewKit(
	writer output.Writer,
	tagLister cli.TagLister,
	tagCreator cli.TagCreator,
	tagPusher cli.TagPusher,
	tagPublisher publish.Publisher) *Kit {
	return &Kit{
		writer:       writer,
		tagLister:    tagLister,
		tagCreator:   tagCreator,
		tagPusher:    tagPusher,
		tagPublisher: tagPublisher,
	}
}

// A Kit contains all the things needed for creating and publishing a new tag.
type Kit struct {
	writer       output.Writer
	tagLister    cli.TagLister
	tagCreator   cli.TagCreator
	tagPusher    cli.TagPusher
	tagPublisher publish.Publisher
}
