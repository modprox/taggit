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

func Test_Parse_ok(t *testing.T) {
	tag, ok := Parse("v1.3.5")
	require.True(t, ok)
	require.Equal(t, Tag{
		Major: 1,
		Minor: 3,
		Patch: 5,
	}, tag)
}

func Test_Parse_big_ok(t *testing.T) {
	tag, ok := Parse("v111.222.333")
	require.True(t, ok)
	require.Equal(t, Tag{
		Major: 111,
		Minor: 222,
		Patch: 333,
	}, tag)
}

func Test_String(t *testing.T) {
	tag := New(1, 2, 3)
	s := tag.String()
	require.Equal(t, "v1.2.3", s)
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
