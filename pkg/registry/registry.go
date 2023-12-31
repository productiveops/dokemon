package registry

import (
	"context"
	"time"

	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/types"
)

func GetImageDigest(imageName string) (string, error) {
	image, err := ParseImage(imageName)
	if err != nil {
		return "", err
	}

	imageString := image.String()
	ref, err := ImageReference(imageString)
	if err != nil {
		return "", err
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30 * time.Second)
	defer cancel()

	digest, err := docker.GetDigest(ctx, &types.SystemContext{}, ref)
	if err != nil {
		return "", err
	}

	return digest.String(), nil
}