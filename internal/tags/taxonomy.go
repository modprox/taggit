package tags

import (
	"sort"

	"gophers.dev/pkgs/semantic"
)

type Taxonomy map[Triple][]semantic.Tag

func (t Taxonomy) Add(tag semantic.Tag) {
	base := tag.Base()
	triple := NewTriple(base.Major, base.Minor, base.Patch)
	t[triple] = append(t[triple], tag)
}

func (t Taxonomy) Sort() {
	for base := range t {
		sort.Sort(
			sort.Reverse(
				semantic.BySemver(t[base]),
			),
		)
	}
}

// Bases return the list of base versions in order.
func (t Taxonomy) Bases() []Triple {
	bases := make([]Triple, 0, len(t))

	for base := range t {
		bases = append(bases, base)
	}

	sort.Slice(bases, func(a, b int) bool {
		A, B := bases[a], bases[b]
		return A.Less(B)
	})

	return bases
}

// Latest returns the most high semver tag, whether it is a base
// tag or a pre-release tag.
func (t Taxonomy) Latest() semantic.Tag {
	bases := t.Bases()
	base := bases[len(bases)-1]
	return t[base][0]
}
