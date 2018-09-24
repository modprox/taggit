// Command taggit provides a convenience wrapper around `git tag` commands.
package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/modprox/mp/pkg/clients/registry"
	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/pkg/netservice"
	"github.com/modprox/taggit/internal/cli"
	"github.com/modprox/taggit/internal/git"
	"github.com/modprox/taggit/internal/publish"
)

const (
	registryEnv = "TAGGIT_REGISTRY_URL"
)

func main() {
	if len(os.Args) != 2 {
		cli.Usage(1)
	}

	command := os.Args[1]
	gitCmd := git.New("git")

	registryURL := os.Getenv(registryEnv)
	publisher, err := newPublisher(registryURL)
	if err != nil {
		die(err)
	}

	tool := cli.NewTool(os.Stdout, gitCmd, publisher)

	tags, err := git.ListTags(gitCmd)
	if err != nil {
		die(err)
	}

	switch command {
	case "help":
		cli.Usage(0)
	case "list":
		err = tool.List(tags)
	case "zero":
		err = tool.Zero(tags)
	case "patch":
		err = tool.Patch(tags)
	case "minor":
		err = tool.Minor(tags)
	case "major":
		err = tool.Major(tags)
	default:
		cli.Usage(1)
	}

	if err != nil {
		die(err)
	}
}

func die(err error) {
	fmt.Fprintf(os.Stderr, "failed: %v\n", err)
	os.Exit(1)
}

func newPublisher(registryURL string) (publish.Publisher, error) {
	if registryURL == "" {
		return publish.Discard(), nil
	}

	parsedURL, err := url.Parse(registryURL)
	if err != nil {
		return nil, errors.Wrap(err, "not a valid registry URL")
	}

	port, err := strconv.Atoi(parsedURL.Port())
	if err != nil {
		port = 0
	}

	scheme := parsedURL.Scheme + "://"
	instances := []netservice.Instance{{
		Address: scheme + parsedURL.Host,
		Port:    port,
	}}

	client := registry.NewClient(registry.Options{
		Instances: instances,
		Timeout:   10 * time.Second,
		Log:       loggy.Discard(),
	})

	modFinder := publish.NewModFinder()

	return publish.New(client, modFinder), nil
}
