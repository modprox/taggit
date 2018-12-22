package updates

import "time"

type Lister interface {
	List() ([]Update, error)
}

type Options struct {
	Timeout    time.Duration
	Executable string // go command
}

func NewLister(options Options) Lister {
	timeout := options.Timeout
	if timeout == 0 {
		timeout = 1 * time.Minute
	}

	executable := options.Executable
	if executable == "" {
		executable = "go"
	}

	return &lister{
		command: newCmd(executable),
		timeout: timeout,
	}
}

type lister struct {
	command cmd
	timeout time.Duration
}

func (l *lister) List() ([]Update, error) {
	s, err := l.command.Run(l.timeout)
	if err != nil {
		return nil, err
	}
	return parseCmd(s)
}
