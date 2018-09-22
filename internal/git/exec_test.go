package git

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Run(t *testing.T) {
	command := New("/bin/echo")

	result, err := command.Run([]string{
		"-n", "a", "b", "c",
	}, 1*time.Second)

	require.NoError(t, err)
	require.Equal(t, "a b c", result)
}

func Test_Run_fail(t *testing.T) {
	command := New("/bin/foobar")

	_, err := command.Run([]string{
		"a", "b", "c",
	}, 1*time.Second)

	require.Error(t, err)
}
