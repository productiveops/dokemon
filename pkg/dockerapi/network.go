package dockerapi

import (
	"context"
	"sort"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func NetworkList(req *DockerNetworkList) (*DockerNetworkListResponse, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	dnetworks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		return nil, err
	}

	networks := make([]Network, len(dnetworks))
	for i, item := range dnetworks {
		networks[i] = Network{
			Id:		item.ID,
			Name: 	item.Name,
			Driver: item.Driver,
			Scope: 	item.Scope,
		}
	}

	sort.Slice(networks, func(i, j int) bool {
		return networks[i].Name < networks[j].Name
	  })

	return &DockerNetworkListResponse{Items: networks}, nil
}

func NetworkRemove(req *DockerNetworkRemove) (error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	err = cli.NetworkRemove(context.Background(), req.Id)
	if err != nil {
		return err
	}

	return nil
}

func NetworksPrune(req *DockerNetworksPrune) (*DockerNetworksPruneResponse, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	report, err := cli.NetworksPrune(context.Background(), filters.NewArgs())
	if err != nil {
		return nil, err
	}

	return &DockerNetworksPruneResponse{NetworksDeleted: report.NetworksDeleted}, nil
}
