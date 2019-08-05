package tags // import "oss.indeed.com/go/taggit/tags"

import (
	"fmt"
	"regexp"
	"strconv"

	"gophers.dev/pkgs/regexplus"
)

var (
	semverRe = regexp.MustCompile(`^v(?P<major>[0-9]+)\.(?P<minor>[0-9]+)\.(?P<patch>[0-9]+)(-(?P<ext>[a-zA-Z0-9_-]+))?$`)
)

func New(major, minor, patch int) Tag {
	return Tag{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

func New2(major, minor, patch int, extension string) Tag {
	return Tag{
		Major:     major,
		Minor:     minor,
		Patch:     patch,
		Extension: extension,
	}
}

func Parse(s string) (Tag, bool) {
	matches := regexplus.FindNamedSubmatches(semverRe, s)

	major, exists := matches["major"]
	if !exists {
		return Tag{}, false
	}

	minor, exists := matches["minor"]
	if !exists {
		return Tag{}, false
	}

	patch, exists := matches["patch"]
	if !exists {
		return Tag{}, false
	}

	extension := matches["ext"]

	return Tag{
		Major:     number(major),
		Minor:     number(minor),
		Patch:     number(patch),
		Extension: extension,
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
	Major     int
	Minor     int
	Patch     int
	Extension string
}

func (t Tag) String() string {
	base := fmt.Sprintf(
		"v%d.%d.%d",
		t.Major,
		t.Minor,
		t.Patch,
	)

	if t.Extension == "" {
		return base
	}

	return base + "-" + t.Extension
}

func (t Tag) Base() Tag {
	return New(t.Major, t.Minor, t.Patch)
}

func (t Tag) IsBase() bool {
	return t.Extension == ""
}

func (t Tag) Less(o Tag) bool {
	if t.Major < o.Major {
		return true
	} else if t.Major > o.Major {
		return false
	} else if extAlessB(t.Extension, o.Extension) {
		return true
	}

	if t.Minor < o.Minor {
		return true
	} else if t.Minor > o.Minor {
		return false
	} else if extAlessB(t.Extension, o.Extension) {
		return true
	}

	if t.Patch < o.Patch {
		return true
	} else if t.Patch > o.Patch {
		return false
	} else if extAlessB(t.Extension, o.Extension) {
		return true
	}

	return false
}

// return true if a's extension precedes b's extension
// normally this is ascibetical, however the empty string
// is a special case that is higher priority than all else
func extAlessB(a, b string) bool {
	if a == "" {
		return false
	}

	if b == "" {
		return true
	}

	return a < b
}

type BySemver []Tag

func (tags BySemver) Len() int      { return len(tags) }
func (tags BySemver) Swap(x, y int) { tags[x], tags[y] = tags[y], tags[x] }
func (tags BySemver) Less(x, y int) bool {
	X := tags[x]
	Y := tags[y]
	return X.Less(Y)
}
