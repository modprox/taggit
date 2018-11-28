package publish

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type ModFinder interface {
	FindModule(root string) (string, error)
}

type modFinder struct{}

func NewModFinder() ModFinder {
	return &modFinder{}
}

func (m *modFinder) FindModule(cwd string) (string, error) {
	return m.lookIn(cwd)
}

func (m *modFinder) lookIn(path string) (string, error) {
	if path == "/" {
		return "", errors.New("no go.mod file found")
	}

	// list the files in path, looking for go.mod file
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		filename := file.Name()
		if filename == "go.mod" {
			fullName := filepath.Join(path, filename)
			return m.readModFile(fullName)
		}
	}

	upOne := filepath.Dir(path)
	return m.lookIn(upOne)
}

// e.g. module file
// module github.com/modprox/taggit

var (
	moduleRe = regexp.MustCompile(`module[\s]+([\S]+)`)
)

func (m *modFinder) readModFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		groups := moduleRe.FindStringSubmatch(line)
		if len(groups) == 2 {
			fullModuleName := groups[1]
			noSuffixName := trimMajorSuffix(fullModuleName)
			return noSuffixName, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", errors.New("go.mod file did not declare module name")
}

var (
	suffixRe = regexp.MustCompile(`(/v[0-9]+)$`)
)

// modules which are at version v2+ must be suffixed with
// their major version (e.g. foo/bar/v2), however the compiler
// does not use that suffix when resolving modules.
func trimMajorSuffix(module string) string {
	return suffixRe.ReplaceAllString(module, "")
}
