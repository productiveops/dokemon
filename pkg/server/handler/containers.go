package handler

import (
	"dokemon/pkg/dockerapi"
	"dokemon/pkg/messages"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func (h *Handler) GetContainerList(c echo.Context) error {
	var err error

	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	req := dockerapi.DockerContainerList{All: true}

	var res *dockerapi.DockerContainerListResponse
	if nodeId == 1 {
		res, err = dockerapi.ContainerList(&req)
	} else {
		res, err = messages.ProcessTaskWithResponse[dockerapi.DockerContainerList, dockerapi.DockerContainerListResponse](uint(nodeId), req, defaultTimeout)
	}

	if err != nil {
		return unprocessableEntity(c, err)
	}	

	return ok(c, res)
}

func (h *Handler) StartContainer(c echo.Context) error {
	var err error

	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	m :=  dockerapi.DockerContainerStart{}
	r := &dockerContainerStartRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	if nodeId == 1 {
		err = dockerapi.ContainerStart(&m)
	} else {
		err = messages.ProcessTask[dockerapi.DockerContainerStart](uint(nodeId), m, defaultTimeout)
	}

	if err != nil {
		return unprocessableEntity(c, err)
	}

	return noContent(c)
}

func (h *Handler) StopContainer(c echo.Context) error {
	var err error

	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	m :=  dockerapi.DockerContainerStop{}
	r := &dockerContainerStopRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	if nodeId == 1 {
		err = dockerapi.ContainerStop(&m)
	} else {
		err = messages.ProcessTask[dockerapi.DockerContainerStop](uint(nodeId), m, defaultTimeout)
	}

	if err != nil {
		return unprocessableEntity(c, err)
	}

	return noContent(c)
}

func (h *Handler) RestartContainer(c echo.Context) error {
	var err error

	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	m :=  dockerapi.DockerContainerRestart{}
	r := &dockerContainerRestartRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	if nodeId == 1 {
		err = dockerapi.ContainerRestart(&m)
	} else {
		err = messages.ProcessTask[dockerapi.DockerContainerRestart](uint(nodeId), m, defaultTimeout)
	}

	if err != nil {
		return unprocessableEntity(c, err)
	}

	return noContent(c)
}

func (h *Handler) RemoveContainer(c echo.Context) error {
	var err error

	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	m :=  dockerapi.DockerContainerRemove{}
	r := &dockerContainerRemoveRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	if nodeId == 1 {
		err = dockerapi.ContainerRemove(&m)
	} else {
		err = messages.ProcessTask[dockerapi.DockerContainerRemove](uint(nodeId), m, defaultTimeout)
	}

	if err != nil {
		return unprocessableEntity(c, err)
	}

	return noContent(c)
}

func (h *Handler) ViewContainerLogs(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	id := c.Param("id")

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error while upgrading from http to websocket")
		return err
	}
	defer ws.Close()

	req := dockerapi.DockerContainerLogs{Id: id}
	if nodeId == 1 {
		err := dockerapi.ContainerLogs(&req, ws)
		if err != nil {
			log.Debug().Err(err).Msg("Error while calling ContainerLogs")
		}
	} else {
		err = messages.ProcessStreamTask[dockerapi.DockerContainerLogs](uint(nodeId), req, ws)
		if err != nil {
			log.Debug().Err(err).Msg("Error while calling ContainerLogs ProcessStreamTask")
		}
	}

	return nil
}

func (h *Handler) OpenContainerTerminal(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	id := c.Param("id")

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	req := dockerapi.DockerContainerTerminal{Id: id}
	if nodeId == 1 {
		err = dockerapi.ContainerTerminal(&req, ws)
		if err != nil {
			log.Debug().Err(err).Msg("Error while calling ContainerTerminal ProcessStreamTask")
			return unprocessableEntity(c, err)
		}
	} else {
		err = messages.ProcessStreamTask[dockerapi.DockerContainerTerminal](uint(nodeId), req, ws)
		if err != nil {
			log.Debug().Err(err).Msg("Error while calling ContainerTerminal ProcessStreamTask")
		}
	}

	return nil
}