package output

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

type Writer interface {
	Directf(format string, args ...interface{})
	Tracef(format string, args ...interface{})
	Writef(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

const (
	tracesEnv = "TAGGIT_ENABLE_TRACES"
)

func enableTraces(env string) bool {
	tracesOn, err := strconv.ParseBool(os.Getenv(env))
	if err != nil {
		return false
	}
	return tracesOn
}

// NewWriter creates a new Writer with the given output sinks. Typically one
// would plug normal into os.Stdout and failure into os.Stderr, but other
// outputs may be provided, for example in use of test cases.
//
// Setting the TAGGIT_ENABLE_TRACES environment variable to 'true' will turn on
// highly verbose trace statements, useful for debugging while doing development
// work.
func NewWriter(normal, failure io.Writer) Writer {
	return &writer{
		normal:  normal,
		failure: failure,
		traces:  enableTraces(tracesEnv),
	}
}

type writer struct {
	normal  io.Writer
	failure io.Writer
	traces  bool
}

func (w *writer) Directf(format string, args ...interface{}) {
	tweaked := format + "\n"
	s := fmt.Sprintf(tweaked, args...)
	w.write(s)
}

func (w *writer) Tracef(format string, args ...interface{}) {
	if w.traces {
		tweaked := "trace: " + format + "\n"
		s := fmt.Sprintf(tweaked, args...)
		w.write(s)
	}
}

func (w *writer) Writef(format string, args ...interface{}) {
	tweaked := "taggit: " + format + "\n"
	s := fmt.Sprintf(tweaked, args...)
	w.write(s)
}

func (w *writer) Errorf(format string, args ...interface{}) {
	tweaked := "taggit: " + format + "\n"
	s := fmt.Sprintf(tweaked, args...)
	w.error(s)
}

func (w *writer) write(s string) {
	_, _ = io.WriteString(w.normal, s)
}

func (w *writer) error(s string) {
	_, _ = io.WriteString(w.failure, s)
}
