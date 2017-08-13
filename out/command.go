package out

import (
	"time"

	"github.com/crdant/cf-route-resource"
)

type Command struct {
	paas PAAS
}

func NewCommand(paas PAAS) *Command {
	return &Command{
		paas: paas,
	}
}

func (command *Command) Run(request Request) (Response, error) {
	err := command.paas.Login(
		request.Source.API,
		request.Source.Username,
		request.Source.Password,
		request.Source.SkipCertCheck,
	)
	if err != nil {
		return Response{}, err
	}

	err = command.paas.Target(
		request.Source.Organization,
		request.Source.Space,
	)
	if err != nil {
		return Response{}, err
	}

	if request.Params.Create != nil {
		for _, r := range request.Params.Create {

			route, err := ParseRoute(r, request.Params.RandomPort)
			if err != nil {
				return Response{}, err
			}
			err = command.paas.CreateRoute(request.Source.Space, route.domain, route.host, route.path, route.port, route.randomPort)
			if err != nil {
				return Response{}, err
			}

		}
	}

	if request.Params.Map != nil {
		for _, r := range request.Params.Map {

			route, err := ParseRoute(r, request.Params.RandomPort)
			if err != nil {
				return Response{}, err
			}

			err = command.paas.MapRoute(request.Params.Application, route.domain, route.host, route.path, route.port)
			if err != nil {
				return Response{}, err
			}

		}
	}

	if request.Params.Unmap != nil {
		for _, r := range request.Params.Unmap {

			route, err := ParseRoute(r, request.Params.RandomPort)
			if err != nil {
				return Response{}, err
			}

			err = command.paas.UnmapRoute(request.Params.Application, route.domain, route.host, route.path, route.port)
			if err != nil {
				return Response{}, err
			}

		}
	}

	return Response{
		Version: resource.Version{
			Timestamp: time.Now(),
		},
		Metadata: []resource.MetadataPair{
			{
				Name:  "organization",
				Value: request.Source.Organization,
			},
			{
				Name:  "space",
				Value: request.Source.Space,
			},
		},
	}, nil
}
