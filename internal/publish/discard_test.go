package publish

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"gophers.dev/pkgs/semantic"

	"oss.indeed.com/go/taggit/internal/cli/output"
)

func Test_Discard(t *testing.T) {
	var (
		normal bytes.Buffer
		bad    bytes.Buffer
	)

	w := output.NewWriter(&normal, &bad)
	p := Discard(w)

	err := p.Publish(semantic.New(1, 2, 3))
	require.NoError(t, err)

	require.Equal(t, "", normal.String())
	require.Equal(t, "", bad.String())
}
