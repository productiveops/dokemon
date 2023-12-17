package dockerapi

import (
	"context"
	"sort"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

func VolumeList(req *DockerVolumeList) (*DockerVolumeListResponse, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	dvolumes, err := cli.VolumeList(context.Background(), volume.ListOptions{})
	if err != nil {
		return nil, err
	}

	volumes := make([]Volume, len(dvolumes.Volumes))
	for i, item := range dvolumes.Volumes {
		volumes[i] = Volume{
			Driver: item.Driver,
			Name: item.Name,
		}
	}

	sort.Slice(volumes, func(i, j int) bool {
		return volumes[i].Name < volumes[j].Name
	  })

	return &DockerVolumeListResponse{Items: volumes}, nil
}

func VolumeRemove(req *DockerVolumeRemove) (error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	err = cli.VolumeRemove(context.Background(), req.Name, true)
	if err != nil {
		return err
	}

	return nil
}

func VolumesPrune(req *DockerVolumesPrune) (*DockerVolumesPruneResponse, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	all := "true"
	if req.All {
		all = "false"
	}
	allFilter := filters.KeyValuePair{Key: "all", Value: all}

	report, err := cli.VolumesPrune(context.Background(), filters.NewArgs(allFilter))
	if err != nil {
		return nil, err
	}

	return &DockerVolumesPruneResponse{VolumesDeleted: report.VolumesDeleted, SpaceReclaimed: report.SpaceReclaimed}, nil
}