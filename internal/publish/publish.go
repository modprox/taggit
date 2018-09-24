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

//go:generate mockery -interface=Publisher -package=publishtest

type Publisher interface {
	Publish(version string) error
}

type publisher struct {
	client    registry.Client
	modFinder ModFinder
}

func New(client registry.Client, modFinder ModFinder) Publisher {
	return &publisher{
		client:    client,
		modFinder: modFinder,
	}
}

// Publish module of name and version to the registry.
func (p *publisher) Publish(version string) error {
	module, err := p.modFinder.FindModule(".")
	if err != nil {
		return err
	}

	wantToAdd := []coordinates.Module{
		{Source: module, Version: version},
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

func (p *discard) Publish(version string) error {
	return nil
}
