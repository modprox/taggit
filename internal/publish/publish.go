package publish

import (
	"bytes"
	"encoding/json"
	"os"

	"gophers.dev/pkgs/semantic"

	"oss.indeed.com/go/modprox/pkg/clients/registry"
	"oss.indeed.com/go/modprox/pkg/coordinates"
	"oss.indeed.com/go/modprox/pkg/netservice"
	"oss.indeed.com/go/taggit/internal/cli/output"
)

const (
	RegistryEnv   = "TAGGIT_REGISTRY_URL"
	mpNewEndpoint = "/v1/registry/sources/new"
)

//go:generate go run github.com/gojuno/minimock/v3/cmd/minimock -g -i Publisher -s _mock.go

type Publisher interface {
	Publish(semantic.Tag) error
}

type publisher struct {
	client    registry.Client
	modFinder ModFinder
	writer    output.Writer
}

func FromEnv(env string, w output.Writer) Publisher {
	registryURL := os.Getenv(env)
	if registryURL == "" {
		return Discard(w)
	}

	client := registry.NewClient(
		registry.Options{
			Instances: []netservice.Instance{{
				Address: registryURL,
			}},
		},
	)

	modFinder := NewModFinder()

	return New(client, modFinder, w)
}

func New(client registry.Client, modFinder ModFinder, w output.Writer) Publisher {
	return &publisher{
		client:    client,
		modFinder: modFinder,
		writer:    w,
	}
}

// Publish module of name and version to the registry.
func (p *publisher) Publish(tag semantic.Tag) error {
	version := tag.String()

	module, err := p.modFinder.FindModule(".")
	if err != nil {
		return err
	}

	wantToAdd := []coordinates.Module{{
		Source:  module,
		Version: version,
	}}

	bs, err := json.Marshal(wantToAdd)
	if err != nil {
		return err
	}
	body := bytes.NewReader(bs)

	var response bytes.Buffer
	if err := p.client.Post(mpNewEndpoint, body, &response); err != nil {
		// maybe read and print body
		return err
	}

	p.writer.Writef("published %s @ %s", module, tag)

	return nil
}
