package dockerapi

import (
	"context"
	"errors"
	"slices"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func ImageList(req *DockerImageList) (*DockerImageListResponse, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	dcontainers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: req.All})
	if err != nil {
		return nil, err
	}

	dimages, err := cli.ImageList(context.Background(), types.ImageListOptions{All: req.All})
	if err != nil {
		return nil, err
	}

	images := make([]Image, len(dimages))
	for i, item := range dimages {
		name := "<none>"
		tag := "<none>"
		inUse := slices.ContainsFunc(dcontainers, func(c types.Container) bool { return c.ImageID == item.ID })
		untagged := len(item.RepoTags) == 0 

		if untagged {
			if len(item.RepoDigests) >= 1 {
				name = strings.Split(item.RepoDigests[0], "@")[0]
			}
		} else {
			nameAndTag := strings.Split(item.RepoTags[0], ":")
			if len(nameAndTag) == 2 {
				name = nameAndTag[0]
				tag = nameAndTag[1]
			}
		}

		images[i] = Image{
			Id:			item.ID,
			Name: 		name,
			Tag: 		tag,
			Size:   	item.Size,
			Dangling: 	untagged && !inUse,
			Created: 	item.Created,
			InUse: 		inUse,
		}
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].Name < images[j].Name
	  })

	return &DockerImageListResponse{Items: images}, nil
}

func ImageRemove(req *DockerImageRemove) (error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	response, err := cli.ImageRemove(context.Background(), req.Id, types.ImageRemoveOptions{Force: req.Force})
	if err != nil {
		return err
	}

	if len(response) == 0 {
		return errors.New("delete unsuccessful")
	}

	return nil
}

func ImagesPrune(req *DockerImagesPrune) (*DockerImagesPruneResponse, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	dangling := "true"
	if req.All {
		dangling = "false"
	}
	danglingFilter := filters.KeyValuePair{Key: "dangling", Value: dangling}

	report, err := cli.ImagesPrune(context.Background(), filters.NewArgs(danglingFilter))
	if err != nil {
		return nil, err
	}
	imagesDeleted := make([]DockerImagesPruneDeletedItem, len(report.ImagesDeleted))
	for i, item := range report.ImagesDeleted {
		imagesDeleted[i] = DockerImagesPruneDeletedItem{Deleted: item.Deleted, Untagged: item.Untagged}
	}

	return &DockerImagesPruneResponse{ImagesDeleted: imagesDeleted, SpaceReclaimed: report.SpaceReclaimed}, nil
}
