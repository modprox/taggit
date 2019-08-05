package cli

import (
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"gophers.dev/pkgs/loggy"

	"oss.indeed.com/go/modprox/pkg/clients/registry"
	"oss.indeed.com/go/modprox/pkg/netservice"
	"oss.indeed.com/go/taggit/internal/git"
	"oss.indeed.com/go/taggit/internal/publish"
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
	tags      Groups
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

	tagsFromGit, err := git.ListTags(gitCmd)
	if err != nil {
		return nil, err
	}

	tagsInGroups, err := Parse(tagsFromGit)
	if err != nil {
		return nil, err
	}

	return &starter{
		args:      arguments,
		publisher: publisher,
		gitCmd:    gitCmd,
		cliTool:   newTool(os.Stdout, gitCmd, publisher),
		tags:      tagsInGroups,
	}, nil
}

func (s *starter) Start() error {
	var err error
	switch s.args.command {
	case "help":
		Usage(0)
	case "list":
		err = s.cliTool.List(s.tags)
	case "zero":
		err = s.cliTool.Zero(s.tags)
	case "patch":
		err = s.cliTool.Patch(s.tags, s.args.extension)
	case "minor":
		err = s.cliTool.Minor(s.tags, s.args.extension)
	case "major":
		err = s.cliTool.Major(s.tags, s.args.extension)
	default:
		Usage(1)
	}

	return err
}

type arguments struct {
	command   string
	extension string
}

func parseArguments(args []string) (arguments, error) {
	if len(args) < 2 || len(args) > 3 {
		return arguments{}, errors.New("wrong number of arguments")
	}

	command := args[1]
	extension := ""
	if len(args) == 3 {
		extension = args[2]
	}

	return arguments{
		command:   command,
		extension: extension,
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
