package tags

import (
	"gophers.dev/pkgs/semantic"
)

func HasPrevious(tax Taxonomy) bool {
	return len(tax) > 0
}

// prev.PR = no, next.PR = no   => bump
// prev.PR = yes, next.PR = no  => no (finalize)
// prev.PR = no, next.PR = yes  => bump
// prev.PR = yes, next.PR = yes => no (continue)

func same(tag semantic.Tag, ext Extensions) semantic.Tag {
	return semantic.New3(
		tag.Major,
		tag.Minor,
		tag.Patch,
		ext.PreRelease,
		ext.BuildMetadta,
	)
}

func IncMajor(previous semantic.Tag, ext Extensions) semantic.Tag {
	if !previous.IsBase() {
		return same(previous, ext)
	}

	base := previous.Base()
	return semantic.New3(
		base.Major+1,
		0,
		0,
		ext.PreRelease,
		ext.BuildMetadta,
	)
}

func IncMinor(previous semantic.Tag, ext Extensions) semantic.Tag {
	if !previous.IsBase() {
		return same(previous, ext)
	}

	base := previous.Base()
	return semantic.New3(
		base.Major,
		base.Minor+1,
		0,
		ext.PreRelease,
		ext.BuildMetadta,
	)
}

func IncPatch(previous semantic.Tag, ext Extensions) semantic.Tag {
	if !previous.IsBase() {
		return same(previous, ext)
	}

	base := previous.Base()
	return semantic.New3(
		base.Major,
		base.Minor,
		base.Patch+1,
		ext.PreRelease,
		ext.BuildMetadta,
	)
}
