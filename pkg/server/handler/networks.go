package handler

import (
	"dokemon/pkg/dockerapi"
	"dokemon/pkg/messages"
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetNetworkList(c echo.Context) error {
	var err error

	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}
	
	req := dockerapi.DockerNetworkList{}

	var res *dockerapi.DockerNetworkListResponse
	if nodeId == 1 {
		res, err = dockerapi.NetworkList(&req)
	} else {
		res, err = messages.ProcessTaskWithResponse[dockerapi.DockerNetworkList, dockerapi.DockerNetworkListResponse](uint(nodeId), req, defaultTimeout)
	}

	if err != nil {
		return unprocessableEntity(c, err)
	}	

	return ok(c, res)
}

func (h *Handler) RemoveNetwork(c echo.Context) error {
	var err error

	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	m := dockerapi.DockerNetworkRemove{}
	r := &dockerNetworkRemoveRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	if nodeId == 1 {
		err = dockerapi.NetworkRemove(&m)
	} else {
		err = messages.ProcessTask[dockerapi.DockerNetworkRemove](uint(nodeId), m, defaultTimeout)
	}

	if err != nil {
		return unprocessableEntity(c, err)
	}

	return noContent(c)
}

func (h *Handler) PruneNetworks(c echo.Context) error {
	var err error

	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	m := dockerapi.DockerNetworksPrune{}

	var res *dockerapi.DockerNetworksPruneResponse
	if nodeId == 1 {
		res, err = dockerapi.NetworksPrune(&m)
	} else {
		res, err = messages.ProcessTaskWithResponse[dockerapi.DockerNetworksPrune, dockerapi.DockerNetworksPruneResponse](uint(nodeId), m, defaultTimeout)
	}

	if err != nil {
		return unprocessableEntity(c, err)
	}

	return ok(c, res)
}
