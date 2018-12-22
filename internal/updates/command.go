package updates

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"
)

//go:generate mockery3 -interface cmd -package updatestest

// cmd returns the command output of running `go list -u -m all`
type cmd interface {
	Run(timeout time.Duration) (string, error)
}

type goCmd struct {
	executable string
}

// newCmd creates a cmd which executes go, given the provided executable.
// If executable is the empty string, the first thing providing 'go' on PATH is
// used.
func newCmd(executable string) cmd {
	return &goCmd{
		executable: executable,
	}
}

func (l *goCmd) Run(timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		l.executable,
		"list", "-m", "-u", "all",
	)
	cmd.Env = os.Environ() // use the tty's environment

	bs, err := cmd.CombinedOutput()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "subprocess failed: %v\n", err.Error())
		_, _ = fmt.Fprintf(os.Stderr, "output: %s\n", string(bs))
		return "", errors.Wrap(err, "go command failed")
	}

	return string(bs), nil
}
