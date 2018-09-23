package publish

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func writeFile(t *testing.T, filename, content string) {
	t.Log("writing file:", filename)
	err := ioutil.WriteFile(filename, []byte(content), 0660)
	require.NoError(t, err)
}

func cleanupFile(t *testing.T, filename string) {
	t.Log("cleanup file:", filename)
	err := os.Remove(filename)
	require.NoError(t, err)
}

func Test_modFinder_readModFile_ok(t *testing.T) {
	writeFile(t, "/tmp/go.mod", "module a/b/c")
	defer cleanupFile(t, "/tmp/go.mod")

	finder := NewModFinder().(*modFinder)
	mod, err := finder.readModFile("/tmp/go.mod")
	require.NoError(t, err)
	require.Equal(t, "a/b/c", mod)
}

func Test_modFinder_readModFile_invalid(t *testing.T) {
	writeFile(t, "/tmp/go.mod", "nothing\nin\nhere")
	defer cleanupFile(t, "/tmp/go.mod")

	finder := NewModFinder().(*modFinder)
	_, err := finder.readModFile("/tmp/go.mod")
	require.Error(t, err)
}

func Test_modFinder_readModFile_missing(t *testing.T) {
	finder := NewModFinder().(*modFinder)
	_, err := finder.readModFile("/tmp/nowhere/go.mod")
	require.Error(t, err)
}

func makeDirs(t *testing.T, dirs string) string {
	root, err := ioutil.TempDir("", "taggit-")
	require.NoError(t, err)

	fullPath := filepath.Join(root, dirs)
	t.Log("creating directories:", fullPath)

	err = os.MkdirAll(fullPath, 0770)
	require.NoError(t, err)

	return root
}

func cleanupDirs(t *testing.T, root string) {
	t.Log("cleanup directories from:", root)
	err := os.RemoveAll(root)
	require.NoError(t, err)
}

func Test_modFinder_FindModule_cwd(t *testing.T) {
	root := makeDirs(t, "foo/bar/baz")
	writeFile(t, filepath.Join(root, "foo/bar/baz", "go.mod"), "module a/b/c")
	defer cleanupDirs(t, root)

	finder := NewModFinder()
	mod, err := finder.FindModule(filepath.Join(root, "foo/bar/baz"))
	require.NoError(t, err)
	require.Equal(t, "a/b/c", mod)
}

func Test_modFinder_FindModule_subdirs(t *testing.T) {
	deep := "foo/bar/baz/one/two/three/four/five"
	root := makeDirs(t, deep)
	writeFile(t, filepath.Join(root, "foo/bar", "go.mod"), "module a/b/c")
	defer cleanupDirs(t, root)

	finder := NewModFinder()
	mod, err := finder.FindModule(filepath.Join(root, deep))
	require.NoError(t, err)
	require.Equal(t, "a/b/c", mod)
}

func Test_modFinder_FindModule_invalid(t *testing.T) {
	root := makeDirs(t, "foo/bar/baz")
	writeFile(t, filepath.Join(root, "foo/bar/baz", "go.mod"), "invalid")
	defer cleanupDirs(t, root)

	finder := NewModFinder()
	_, err := finder.FindModule(filepath.Join(root, "foo/bar/baz"))
	require.Error(t, err)
}

func Test_modFinder_FindModule_noFile(t *testing.T) {
	deep := "foo/bar/baz/one/two/three/four/five"
	root := makeDirs(t, deep)
	defer cleanupDirs(t, root)

	finder := NewModFinder()
	_, err := finder.FindModule(filepath.Join(root, deep))
	require.Error(t, err)
}
