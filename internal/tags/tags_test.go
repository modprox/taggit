package tags

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Triple_Major(t *testing.T) {
	triple := NewTriple(3, 4, 5)
	major := triple.Major()
	require.Equal(t, 3, major)
}

func Test_Triple_Minor(t *testing.T) {
	triple := NewTriple(3, 4, 5)
	minor := triple.Minor()
	require.Equal(t, 4, minor)
}

func Test_Triple_Patch(t *testing.T) {
	triple := NewTriple(3, 4, 5)
	patch := triple.Patch()
	require.Equal(t, 5, patch)
}

func Test_Triple_String(t *testing.T) {
	triple := NewTriple(3, 4, 5)
	s := triple.String()
	require.Equal(t, "v3.4.5", s)
}

func Test_Triple_Less(t *testing.T) {
	try := func(a, b Triple, expLess bool) {
		result := a.Less(b)
		require.Equal(t, expLess, result)
	}

	try(
		NewTriple(0, 0, 1),
		NewTriple(0, 0, 2),
		true,
	)

	try(
		NewTriple(0, 0, 2),
		NewTriple(0, 0, 1),
		false,
	)

	try(
		NewTriple(0, 1, 0),
		NewTriple(0, 2, 0),
		true,
	)

	try(
		NewTriple(0, 2, 0),
		NewTriple(0, 1, 0),
		false,
	)

	try(
		NewTriple(1, 0, 0),
		NewTriple(2, 0, 0),
		true,
	)

	try(
		NewTriple(2, 0, 0),
		NewTriple(1, 0, 0),
		false,
	)
}
