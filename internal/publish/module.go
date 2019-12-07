package publish

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"gophers.dev/pkgs/ignore"
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i ModFinder -s _mock.go

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
// module oss.indeed.com/go/taggit

var (
	moduleRe = regexp.MustCompile(`module[\s]+([\S]+)`)
)

func (m *modFinder) readModFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer ignore.Close(f)

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
