package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"oss.indeed.com/go/taggit/tags"
)

const sampleTags = `
deploy/2017-03-17
deploy/2017-05-12
deploy/2017-10-12_09-05-04
deploy/2017-10-12_09-53-03
deploy/2017-10-12_10-42-32
v0.0.0
v0.0.1
v0.0.1-alpha
v0.0.1-alpha2
v0.0.5
v0.0.6
v1.0.0
v1.0.0-rc1
v1.1.0
v1.1.1
0.0.3
`

func Test_Groups_Parse(t *testing.T) {
	groups, err := Parse(sampleTags)
	require.NoError(t, err)

	exp := Groups{
		tags.New(0, 0, 0): {
			tags.New(0, 0, 0),
		},
		tags.New(0, 0, 1): {
			tags.New(0, 0, 1),
			tags.New2(0, 0, 1, "alpha2"),
			tags.New2(0, 0, 1, "alpha"),
		},
		tags.New(0, 0, 5): {
			tags.New(0, 0, 5),
		},
		tags.New(0, 0, 6): {
			tags.New(0, 0, 6),
		},
		tags.New(1, 0, 0): {
			tags.New(1, 0, 0),
			tags.New2(1, 0, 0, "rc1"),
		},
		tags.New(1, 1, 0): {
			tags.New(1, 1, 0),
		},
		tags.New(1, 1, 1): {
			tags.New(1, 1, 1),
		},
	}

	compareGroups(t, exp, groups)
}

func compareGroups(t *testing.T, exp, groups Groups) {
	if len(exp) != len(groups) {
		t.Fatalf("exp of size %d, got size %d", len(exp), len(groups))
	}

	for expBase, expList := range exp {
		gList, exists := groups[expBase]
		if !exists {
			t.Fatalf("exp had base %v, but groups does not", expBase)
		}

		if len(expList) != len(gList) {
			t.Fatalf(
				"exp list for base %v size %d, got %d",
				expBase,
				len(expList),
				len(gList),
			)
		}

		for i := 0; i < len(expList); i++ {
			if expList[i] != gList[i] {
				t.Fatalf(
					"exp at base %v, index %d: %v, got: %v",
					expBase, i, expList[i], gList[i],
				)
			}
		}
	}
}

func Test_Groups_Bases(t *testing.T) {
	groups, err := Parse(sampleTags)
	require.NoError(t, err)

	bases := groups.Bases()
	require.Equal(t, []tags.Tag{
		tags.New(0, 0, 0),
		tags.New(0, 0, 1),
		tags.New(0, 0, 5),
		tags.New(0, 0, 6),
		tags.New(1, 0, 0),
		tags.New(1, 1, 0),
		tags.New(1, 1, 1),
	}, bases)
}
