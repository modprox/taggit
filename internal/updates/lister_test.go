package updates

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Lister_ok(t *testing.T) {
	goCmd := "../../hack/go-list.sh"

	l := NewLister(Options{
		Timeout:    1 * time.Second,
		Executable: goCmd,
	})

	updates, err := l.List()
	require.NoError(t, err)

	fmt.Println("updates:", updates)
}
