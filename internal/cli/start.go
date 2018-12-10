package cli

import (
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/modprox/mp/pkg/clients/registry"
	"github.com/modprox/mp/pkg/loggy"
	"github.com/modprox/mp/pkg/netservice"
	"github.com/modprox/taggit/internal/git"
	"github.com/modprox/taggit/internal/publish"
	"github.com/modprox/taggit/tags"
)

const (
	registryEnv = "TAGGIT_REGISTRY_URL"
)

type Starter interface {
	Start() error
}

type starter struct {
	args      arguments
	publisher publish.Publisher
	gitCmd    git.Cmd
	cliTool   Tool
	tags      []tags.Tag
}

func New(args []string) (Starter, error) {
	arguments, err := parseArguments(args)
	if err != nil {
		return nil, err
	}

	registryURL := os.Getenv(registryEnv)
	publisher, err := newPublisher(registryURL)
	if err != nil {
		return nil, err
	}

	gitCmd := git.New("git")

	repoTags, err := git.ListTags(gitCmd)
	if err != nil {
		return nil, err
	}

	return &starter{
		args:      arguments,
		publisher: publisher,
		gitCmd:    gitCmd,
		cliTool:   newTool(os.Stdout, gitCmd, publisher),
		tags:      repoTags,
	}, nil
}

func (s *starter) Start() error {
	var err error
	switch s.args.firstCmd {
	case "help":
		Usage(0)
	case "list":
		err = s.cliTool.List(s.tags)
	case "zero":
		err = s.cliTool.Zero(s.tags)
	case "patch":
		err = s.cliTool.Patch(s.tags) // , s.args.secondCmd)
	case "minor":
		err = s.cliTool.Minor(s.tags) // , s.args.secondCmd)
	case "major":
		err = s.cliTool.Major(s.tags) // , s.args.secondCmd)
	default:
		Usage(1)
	}

	return err
}

type arguments struct {
	firstCmd  string
	secondCmd string
}

func parseArguments(args []string) (arguments, error) {
	if len(args) < 2 || len(args) > 3 {
		return arguments{}, errors.New("wrong number of arguments")
	}

	return arguments{
		firstCmd:  args[1],
		secondCmd: args[2],
	}, nil
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
