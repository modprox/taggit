package publish

import (
	"gophers.dev/pkgs/semantic"

	"oss.indeed.com/go/taggit/internal/cli/output"
)

func Discard(w output.Writer) Publisher {
	return &discard{
		writer: w,
	}
}

type discard struct {
	writer output.Writer
}

func (d *discard) Publish(tag semantic.Tag) error {
	d.writer.Tracef("no publisher configured, not publishing tag %s", tag)
	return nil
}
