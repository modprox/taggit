package updates

import "time"

type Lister interface {
	List() ([]Update, error)
}

type Options struct {
	Timeout time.Duration
}

func NewLister(options Options) Lister {
	timeout := options.Timeout
	if timeout == 0 {
		timeout = 1 * time.Minute
	}
	return &lister{
		command: newCmd("go"), // use "go" on $PATH
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
