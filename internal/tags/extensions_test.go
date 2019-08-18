package tags

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ExtractExtensions(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	fs.String("meta", "bm1", "set metadata")

	err := fs.Set("meta", "abc123")
	require.NoError(t, err)

	ext := ExtractExtensions(fs)
	require.Equal(t, "abc123", ext.BuildMetadta)
}
