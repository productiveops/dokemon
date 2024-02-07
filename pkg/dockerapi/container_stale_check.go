package dockerapi

import (
	"context"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/productiveops/dokemon/pkg/registry"
	"github.com/rs/zerolog/log"
)

var containerStaleStatus map[string]string
const (
	StaleStatusProcessing = "processing"
	StaleStatusYes = "yes"
	StaleStatusNo = "no"
	StaleStatusError = "error"
)

func isContainerImageStale(imageAndTag string, imageId string, cli *client.Client) (bool, error) {
	latestDigest, err := registry.GetImageDigest(imageAndTag)
	if err != nil {
		return false, err
	}

	imageInspect, _, err := cli.ImageInspectWithRaw(context.Background(), imageId)
	if err != nil {
		return	false, err
	}

	currentDigest := imageInspect.RepoDigests[0]
	if strings.Contains(currentDigest, "@") {
		currentDigest = strings.Split(currentDigest, "@")[1]
	}

	isStale := currentDigest != latestDigest
	return isStale, nil
}

func ContainerScheduleRefreshStaleStatus() {
	for {
		log.Info().Msg("Refreshing container stale status")
		ContainerRefreshStaleStatus()
		time.Sleep(24 * time.Hour)
	}
}

func ContainerRefreshStaleStatus() error {
	if containerStaleStatus == nil {
		containerStaleStatus = make(map[string]string)
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	dcontainers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return err
	}

	for _, c := range dcontainers {
		_, ok := containerStaleStatus[c.ID]
		if !ok {
			containerStaleStatus[c.ID] = StaleStatusProcessing
		}
	}

	for _, c := range dcontainers {
		image := strings.Split(c.Image, "@")[0]
		stale := StaleStatusProcessing
		isStale, err := isContainerImageStale(image, c.ImageID, cli)
		if err != nil {
			stale = StaleStatusError
			log.Error().Err(err).Str("containerId", c.ID).Str("image", image).Msg("Error while checking if container is stale")
		} else {
			if isStale {
				stale = StaleStatusYes
			} else {
				stale = StaleStatusNo
			}	
		}
		containerStaleStatus[c.ID] = stale
		containerStaleStatus[c.ID[:12]] = stale	// docker compose -p PROJECT ps --format json returns 12 chars of ID. So need this.
	}

	return nil
}