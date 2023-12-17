package handler

import (
	"dokemon/pkg/dockerapi"
	"dokemon/pkg/messages"
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetImageList(c echo.Context) error {
	var err error

	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}
	
	req := dockerapi.DockerImageList{All: true}

	var res *dockerapi.DockerImageListResponse
	if nodeId == 1 {
		res, err = dockerapi.ImageList(&req)
	} else {
		res, err = messages.ProcessTaskWithResponse[dockerapi.DockerImageList, dockerapi.DockerImageListResponse](uint(nodeId), req, defaultTimeout)
	}

	if err != nil {
		return unprocessableEntity(c, err)
	}	

	return ok(c, res)
}

func (h *Handler) RemoveImage(c echo.Context) error {
	var err error

	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	m := dockerapi.DockerImageRemove{}
	r := &dockerImageRemoveRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	if nodeId == 1 {
		err = dockerapi.ImageRemove(&m)
	} else {
		err = messages.ProcessTask[dockerapi.DockerImageRemove](uint(nodeId), m, defaultTimeout)
	}

	if err != nil {
		return unprocessableEntity(c, err)
	}

	return noContent(c)
}

func (h *Handler) PruneImages(c echo.Context) error {
	var err error

	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	m := dockerapi.DockerImagesPrune{}
	r := &dockerImagesPruneRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	var res *dockerapi.DockerImagesPruneResponse
	if nodeId == 1 {
		res, err = dockerapi.ImagesPrune(&m)
	} else {
		res, err = messages.ProcessTaskWithResponse[dockerapi.DockerImagesPrune, dockerapi.DockerImagesPruneResponse](uint(nodeId), m, defaultTimeout)
	}

	if err != nil {
		return unprocessableEntity(c, err)
	}

	return ok(c, res)
}
