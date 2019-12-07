package commands

import "errors"

var (
	ErrNoPreviousTags = errors.New("no previous tags")
	ErrPreviousTags   = errors.New("previous tags exist")
)
