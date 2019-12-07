package tags

import (
	"gophers.dev/pkgs/semantic"
)

func HasPrevious(tax Taxonomy) bool {
	return len(tax) > 0
}

func IncMajor(tag semantic.Tag) semantic.Tag {
	base := tag.Base()
	return semantic.New(
		base.Major+1,
		0,
		0,
	)
}

func IncMinor(tag semantic.Tag) semantic.Tag {
	base := tag.Base()
	return semantic.New(
		base.Major,
		base.Minor+1,
		0,
	)
}

func IncPatch(tag semantic.Tag) semantic.Tag {
	base := tag.Base()
	return semantic.New(
		base.Major,
		base.Minor,
		base.Patch+1,
	)
}
