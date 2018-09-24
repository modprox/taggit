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
			fullname := filepath.Join(path, filename)
			return m.readModFile(fullname)
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
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		groups := moduleRe.FindStringSubmatch(line)
		if len(groups) == 2 {
			return groups[1], nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", errors.New("go.mod file did not declare module name")
}