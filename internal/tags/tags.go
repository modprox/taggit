// Package tags contains utility functions around organizing tags.
package tags

import (
	"gophers.dev/pkgs/semantic"
)

// A Triple represents the core Major, Minor, and Patch levels of a tag. This is
// not enough information to fully define a tag (use semantic.Tag for that), as
// tags may also have associated metadata.
type Triple struct {
	major int
	minor int
	patch int
}

func NewTriple(major, minor, patch int) Triple {
	return Triple{
		major: major,
		minor: minor,
		patch: patch,
	}
}

func (t Triple) Major() int {
	return t.major
}

func (t Triple) Minor() int {
	return t.minor
}

func (t Triple) Patch() int {
	return t.patch
}

func (t Triple) String() string {
	return semantic.New(
		t.Major(),
		t.Minor(),
		t.Patch(),
	).String()
}

func (t Triple) Less(o Triple) bool {
	// first compare major versions
	if t.Major() < o.Major() {
		return true
	} else if t.Major() > o.Major() {
		return false
	}

	// if major are equal, compare minor versions
	if t.Minor() < o.Minor() {
		return true
	} else if t.Minor() > o.Minor() {
		return false
	}

	// if minor are equal, compare patch versions
	return t.Patch() < o.Patch()
}
