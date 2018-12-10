package tags

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	tag := New(1, 2, 3)
	require.Equal(t, Tag{
		Major: 1,
		Minor: 2,
		Patch: 3,
	}, tag)
}

func Test_Extended(t *testing.T) {
	tag := New2(1, 2, 3, "alpha1")
	require.Equal(t, Tag{
		Major:     1,
		Minor:     2,
		Patch:     3,
		Extension: "alpha1",
	}, tag)
}

func Test_Parse(t *testing.T) {
	try := func(s string, exp Tag, expOK bool) {
		result, ok := Parse(s)
		require.Equal(t, expOK, ok)
		require.Equal(t, exp, result)
	}

	try("v1.3.5", New(1, 3, 5), true)
	try("v111.222.333", New(111, 222, 333), true)
	try("v1.2.3-alpha", New2(1, 2, 3, "alpha"), true)
	try("v1.2.3-alpha2", New2(1, 2, 3, "alpha2"), true)
	try("1.2.3", Tag{}, false)       // missing v
	try("v1.2.3_beta", Tag{}, false) // dash required for extension
}

func Test_String(t *testing.T) {
	tag := New(1, 2, 3)
	s := tag.String()
	require.Equal(t, "v1.2.3", s)
}

func Test_String_extension(t *testing.T) {
	tag := New2(1, 2, 3, "alpha2")
	s := tag.String()
	require.Equal(t, "v1.2.3-alpha2", s)
}

func Test_Sort_BySemver(t *testing.T) {
	list := []Tag{
		New(3, 1, 2),
		New(3, 3, 1),
		New(1, 3, 2),
		New(2, 1, 1),
		New(1, 6, 2),
		New(3, 3, 3),
		New(2, 4, 1),
		New(1, 8, 2),
		New(1, 7, 0),
	}
	sort.Sort(sort.Reverse(BySemver(list)))
	require.Equal(t, []Tag{
		New(3, 3, 3),
		New(3, 3, 1),
		New(3, 1, 2),
		New(2, 4, 1),
		New(2, 1, 1),
		New(1, 8, 2),
		New(1, 7, 0),
		New(1, 6, 2),
		New(1, 3, 2),
	}, list)
}

func Test_Sort_BySemver_Extended_sameV(t *testing.T) {
	list := []Tag{
		New2(1, 1, 1, "bbb"),
		New2(1, 1, 1, "ddd"),
		New(1, 1, 1),
		New2(1, 1, 1, "aaa"),
		New2(1, 1, 1, "ccc"),
	}

	sort.Sort(sort.Reverse(BySemver(list)))

	require.Equal(t, []Tag{
		New(1, 1, 1),
		New2(1, 1, 1, "ddd"),
		New2(1, 1, 1, "ccc"),
		New2(1, 1, 1, "bbb"),
		New2(1, 1, 1, "aaa"),
	}, list)
}
