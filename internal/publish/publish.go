package publish

import (
	"bytes"
	"encoding/json"

	"github.com/modprox/mp/pkg/clients/registry"
	"github.com/modprox/mp/pkg/coordinates"
)

const (
	modproxNewAPI = "/v1/registry/sources/new"
)

type Publisher interface {
	Publish(name, version string) error
}

type publisher struct {
	client registry.Client
}

func New(client registry.Client) Publisher {
	return &publisher{
		client: client,
	}
}

// Publish module of name and version to the registry.
func (p *publisher) Publish(name, version string) error {
	wantToAdd := []coordinates.Module{
		{Source: name, Version: version},
	}

	bs, err := json.Marshal(wantToAdd)
	if err != nil {
		return err
	}
	body := bytes.NewReader(bs)

	var response bytes.Buffer
	if err := p.client.Post(modproxNewAPI, body, &response); err != nil {
		// maybe read and print body
		return err
	}
	return nil
}

func Discard() Publisher {
	return &discard{}
}

type discard struct{}

func (p *discard) Publish(name, version string) error {
	return nil
}
