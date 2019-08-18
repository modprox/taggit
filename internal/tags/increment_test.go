package tags

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gophers.dev/pkgs/semantic"
)

func Test_HasPrevious(t *testing.T) {
	tax := Taxonomy{
		NewTriple(1, 2, 3): []semantic.Tag{
			semantic.New(1, 2, 3),
		},
	}
	hasPrevious := HasPrevious(tax)
	require.True(t, hasPrevious)
}

func Test_HasPrevious_empty(t *testing.T) {
	tax := Taxonomy(nil)
	hasPrevious := HasPrevious(tax)
	require.False(t, hasPrevious)
}

func Test_IncMajor(t *testing.T) {
	previous := semantic.New3(1, 2, 3, "rc1", "bm1")
	next := IncMajor(previous)
	require.Equal(t, semantic.New(2, 0, 0), next)
}

func Test_IncMinor(t *testing.T) {
	previous := semantic.New3(1, 2, 3, "rc1", "bm1")
	next := IncMinor(previous)
	require.Equal(t, semantic.New(1, 3, 0), next)
}

func Test_IncPatch(t *testing.T) {
	previous := semantic.New3(1, 3, 3, "rc1", "bm1")
	next := IncPatch(previous)
	require.Equal(t, semantic.New(1, 3, 4), next)
}
