package publish

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"oss.indeed.com/go/modprox/pkg/clients/registry"
	"oss.indeed.com/go/taggit/internal/cli/output"

	"gophers.dev/pkgs/semantic"
)

const (
	egModFile = `module oss.indeed.com/go/taggit`
)

func Test_Publisher_Publish(t *testing.T) {
	client := registry.NewClientMock(t)
	defer client.MinimockFinish()

	modFinder := NewModFinderMock(t)
	defer modFinder.MinimockFinish()

	var a, b bytes.Buffer
	w := output.NewWriter(&a, &b)

	modFinder.FindModuleMock.Expect(".").Return(
		egModFile, nil,
	)

	// client.PostMock.Expect
	client.PostMock.Return(nil)

	p := New(client, modFinder, w)
	err := p.Publish(semantic.New(1, 2, 3))
	require.NoError(t, err)
	require.Equal(t, "taggit: published module oss.indeed.com/go/taggit @ v1.2.3\n", a.String())
	require.Equal(t, "", b.String())
}
