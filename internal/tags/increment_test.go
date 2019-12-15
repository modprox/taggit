package tags

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gophers.dev/pkgs/semantic"
)

func Test_HasPrevious(t *testing.T) {
	t.Parallel()

	try := func(t *testing.T, taxonomy Taxonomy, exp bool) {
		result := HasPrevious(taxonomy)
		require.Equal(t, exp, result)
	}

	t.Run("one existing tag", func(t *testing.T) {
		try(t, Taxonomy{
			NewTriple(1, 2, 3): []semantic.Tag{
				semantic.New(1, 2, 3),
			}},
			true,
		)
	})

	t.Run("no existing tag", func(t *testing.T) {
		try(t, Taxonomy(nil), false)
	})
}

func ext(pr, bm string) Extensions {
	return Extensions{PreRelease: pr, BuildMetadta: bm}
}

func Test_IncMajor(t *testing.T) {
	t.Parallel()

	try := func(t *testing.T, ext Extensions, previous, exp semantic.Tag) {
		next := IncMajor(previous, ext)
		require.Equal(t, exp, next)
	}

	// test the 2x2 matrix of pre-release presence

	t.Run("prevPR.none nextPR.none", func(t *testing.T) {
		// previous tag is final, gets incremented
		try(t,
			ext("", ""),
			semantic.New(1, 2, 3),
			semantic.New(2, 0, 0),
		)
	})

	t.Run("prevPR.yes, nextPR.none", func(t *testing.T) {
		// previous tag is pre-release, gets finalized
		try(t,
			ext("", ""),
			semantic.New2(1, 2, 3, "rc1"),
			semantic.New(1, 2, 3),
		)
	})

	t.Run("prevPR.none, nextPR.yes", func(t *testing.T) {
		// previous tag is final, gets incremented with pre-release
		try(t,
			ext("alpha1", ""),
			semantic.New(1, 2, 3),
			semantic.New2(2, 0, 0, "alpha1"),
		)
	})

	t.Run("prevPR.yes, nextPR.yes", func(t *testing.T) {
		// previous tag is pre-release, next is also pre-release
		try(t,
			ext("beta1", ""),
			semantic.New2(1, 2, 3, "alpha1"),
			semantic.New2(1, 2, 3, "beta1"),
		)
	})

	// test that build metadata is set on next tag

	t.Run("next has build metadata", func(t *testing.T) {
		try(t,
			ext("", "foo"),
			semantic.New(1, 2, 3),
			semantic.New4(2, 0, 0, "foo"),
		)
	})

	t.Run("next overwrites build metadata", func(t *testing.T) {
		try(t,
			ext("rc2", "bar"),
			semantic.New3(1, 2, 3, "rc1", "foo"),
			semantic.New3(1, 2, 3, "rc2", "bar"),
		)
	})
}

func Test_IncMinor(t *testing.T) {
	t.Parallel()

	try := func(t *testing.T, ext Extensions, previous, exp semantic.Tag) {
		next := IncMinor(previous, ext)
		require.Equal(t, exp, next)
	}

	// test the 2x2 matrix of pre-release presence

	t.Run("prevPR.none nextPR.none", func(t *testing.T) {
		// previous tag is final, gets incremented
		try(t,
			ext("", ""),
			semantic.New(1, 2, 3),
			semantic.New(1, 3, 0),
		)
	})

	t.Run("prevPR.yes nextPR.none", func(t *testing.T) {
		// previous tag is pre-release, gets finalized
		try(t,
			ext("", ""),
			semantic.New2(1, 2, 3, "alpha1"),
			semantic.New(1, 2, 3),
		)
	})

	t.Run("prevPR.no nextPR.yes", func(t *testing.T) {
		// previous tag is final, gets incremented with pre-release
		try(t,
			ext("beta1", ""),
			semantic.New(1, 2, 3),
			semantic.New2(1, 3, 0, "beta1"),
		)
	})

	t.Run("prevPR.yes nexPR.yes", func(t *testing.T) {
		// previous tag is pre-release, next is also pre-release
		try(t,
			ext("rc1", ""),
			semantic.New2(1, 2, 3, "alpha1"),
			semantic.New2(1, 2, 3, "rc1"),
		)
	})

	// test that build metadata is set on next tag

	t.Run("next has build metadata", func(t *testing.T) {
		try(t,
			ext("", "foo"),
			semantic.New(1, 2, 3),
			semantic.New4(1, 3, 0, "foo"),
		)
	})

	t.Run("next overwrites build metadata", func(t *testing.T) {
		try(t,
			ext("rc2", "bar"),
			semantic.New3(1, 2, 3, "rc1", "foo"),
			semantic.New3(1, 2, 3, "rc2", "bar"),
		)
	})
}

func Test_IncPatch(t *testing.T) {
	t.Parallel()

	try := func(t *testing.T, ext Extensions, previous, exp semantic.Tag) {
		next := IncPatch(previous, ext)
		require.Equal(t, exp, next)
	}

	// test the 2x2 matrix of pre-release presence

	t.Run("prevPR.none nextPR.none", func(t *testing.T) {
		// previous tag is final, gets incremented
		try(t,
			ext("", ""),
			semantic.New(1, 2, 3),
			semantic.New(1, 2, 4),
		)
	})

	t.Run("prevPR.yes nextPR.none", func(t *testing.T) {
		// previous tag is pre-release, gets finalized
		try(t,
			ext("", ""),
			semantic.New2(1, 2, 3, "rc1"),
			semantic.New(1, 2, 3),
		)
	})

	t.Run("prevPR.none nextPR.yes", func(t *testing.T) {
		// previous tag is final, gets incremented with pre-release
		try(t,
			ext("alpha1", ""),
			semantic.New(1, 2, 3),
			semantic.New2(1, 2, 4, "alpha1"),
		)
	})

	t.Run("prevPR.yes nextPR.yes", func(t *testing.T) {
		// previous tag is pre-release, next is also pre-release
		try(t, ext("beta1", ""),
			semantic.New2(1, 2, 3, "alpha1"),
			semantic.New2(1, 2, 3, "beta1"),
		)
	})

	// test that build metadata is set on next tag

	t.Run("next has build metadata", func(t *testing.T) {
		try(t,
			ext("", "foo"),
			semantic.New(1, 2, 3),
			semantic.New4(1, 2, 4, "foo"),
		)
	})

	t.Run("next overwrites build metadata", func(t *testing.T) {
		try(t,
			ext("rc2", "bar"),
			semantic.New3(1, 2, 3, "rc1", "foo"),
			semantic.New3(1, 2, 3, "rc2", "bar"),
		)
	})
}
