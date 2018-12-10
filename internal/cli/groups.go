package cli

import (
	"bufio"
	"sort"
	"strings"

	"github.com/modprox/taggit/tags"
)

// A Group of tags is all of the tags associated with
// a core semver tag version. The key must be a core tag (no extension).
// The list of tags is each tag (maybe including core tag) at the same
// version as the core tag.
type Groups map[tags.Tag][]tags.Tag

func (g Groups) add(tag tags.Tag) {
	base := tag.Base()
	g[base] = append(g[base], tag)
}

func (g Groups) sort() {
	for base := range g {
		sort.Sort(
			sort.Reverse(
				tags.BySemver(g[base]),
			),
		)
	}
}

// Bases return the list of base versions in order.
func (g Groups) Bases() []tags.Tag {
	bases := make([]tags.Tag, 0, len(g))
	// hello world
	for base := range g {
		bases = append(bases, base)
	}

	sort.Sort(tags.BySemver(bases))

	return bases
}

// Latest returns the most high semver tag, whether it is a base
// tag or a pre-release tag.
func (g Groups) Latest() tags.Tag {
	bases := g.Bases()
	base := bases[len(bases)-1]
	return g[base][0]
}

// Parse the output of `git tag --list` into a set of tags grouped
// such that all extensions are associated with a core version.
func Parse(s string) (Groups, error) {
	groups := make(Groups)
	scanner := bufio.NewScanner(strings.NewReader(s))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if tag, ok := tags.Parse(line); ok {
			groups.add(tag)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	groups.sort()
	return groups, nil
}
