package tags

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	semverRe = regexp.MustCompile(`^v([0-9]+)\.([0-9]+)\.([0-9])$`)
)

func New(major, minor, patch int) Tag {
	return Tag{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

func Parse(s string) (Tag, bool) {
	groups := semverRe.FindStringSubmatch(s)
	if len(groups) != 4 {
		return Tag{}, false
	}

	return Tag{
		Major: number(groups[1]),
		Minor: number(groups[2]),
		Patch: number(groups[3]),
	}, true
}

func number(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic("bug in our tag regexp")
	}
	return n
}

type Tag struct {
	Major int
	Minor int
	Patch int
}

func (t Tag) String() string {
	return fmt.Sprintf(
		"v%d.%d.%d",
		t.Major,
		t.Minor,
		t.Patch,
	)
}

func (t Tag) Less(o Tag) bool {
	if t.Major < o.Major {
		return true
	} else if t.Major > o.Major {
		return false
	}

	if t.Minor < o.Minor {
		return true
	} else if t.Minor > o.Minor {
		return false
	}

	return t.Patch < o.Patch
}

var (
	ZeroValue = Tag{
		Major: 0,
		Minor: 0,
		Patch: 0,
	}
)

type BySemver []Tag

func (tags BySemver) Len() int      { return len(tags) }
func (tags BySemver) Swap(x, y int) { tags[x], tags[y] = tags[y], tags[x] }
func (tags BySemver) Less(x, y int) bool {
	X := tags[x]
	Y := tags[y]
	return X.Less(Y)
}
