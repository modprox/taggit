package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"oss.indeed.com/go/taggit/internal/cli"
	"oss.indeed.com/go/taggit/internal/cli/output"
	"oss.indeed.com/go/taggit/internal/publish"
)

type mocks struct {
	stdout       *bytes.Buffer
	stderr       *bytes.Buffer
	writer       output.Writer
	tagLister    *cli.TagListerMock
	tagCreator   *cli.TagCreatorMock
	tagPublisher *publish.PublisherMock
}

func newMocks(t *testing.T) mocks {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	return mocks{
		stdout:       &stdout,
		stderr:       &stderr,
		writer:       output.NewWriter(&stdout, &stderr),
		tagLister:    cli.NewTagListerMock(t),
		tagCreator:   cli.NewTagCreatorMock(t),
		tagPublisher: publish.NewPublisherMock(t),
	}
}

func (m mocks) assertions(t *testing.T) {
	m.tagLister.MinimockFinish()
	m.tagCreator.MinimockFinish()
	m.tagPublisher.MinimockFinish()
}

func Test_NewKit(t *testing.T) {
	mocks := newMocks(t)
	kit := NewKit(mocks.writer, mocks.tagLister, mocks.tagCreator, mocks.tagPublisher)

	r := require.New(t)
	r.NotNil(t, kit.writer)
	r.NotNil(t, kit.tagLister)
	r.NotNil(t, kit.tagCreator)
	r.NotNil(t, kit.tagPublisher)
}
