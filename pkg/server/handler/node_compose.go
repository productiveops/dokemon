package handler

import (
	"dokemon/pkg/dockerapi"
	"dokemon/pkg/messages"
	"dokemon/pkg/server/model"
	"dokemon/pkg/server/store"
	"errors"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func (h *Handler) GetNodeComposeProjectList(c echo.Context) error {
	var err error

	p, err := strconv.Atoi(c.QueryParam("p"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("p"))
	}

	if p < 1 {
		return unprocessableEntity(c, queryGte1ExpectedError("p"))
	}

	s, err := strconv.Atoi(c.QueryParam("s"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("s"))
	}

	if s < 1 {
		return unprocessableEntity(c, queryGte1ExpectedError("s"))
	}

	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	rows, totalRows, err := h.nodeComposeProjectStore.GetList(uint(nodeId), uint(p), uint(s))
	if err != nil {
		panic(err)
	}

	req := dockerapi.DockerComposeList{}

	var res *dockerapi.DockerComposeListResponse
	if nodeId == 1 {
		res, err = dockerapi.ComposeList(&req)
	} else {
		res, err = messages.ProcessTaskWithResponse[dockerapi.DockerComposeList, dockerapi.DockerComposeListResponse](uint(nodeId), req, defaultTimeout)
	}

	if err != nil {
		return unprocessableEntity(c, err)
	}
	
	return ok(c, newPageResponse(newNodeComposeProjectItemList(rows, res.Items), uint(p), uint(s), uint(totalRows)))
}

func (h *Handler) GetNodeComposeProject(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, errors.New("id should be an integer"))
	}

	ncp, err := h.nodeComposeProjectStore.GetById(uint(nodeId), uint(id))
	if err != nil {
		return unprocessableEntity(c, errors.New("Project not found"))
	}

	req := dockerapi.DockerComposeGet{ProjectName: ncp.ProjectName}

	var res *dockerapi.ComposeItem
	if nodeId == 1 {
		res, err = dockerapi.ComposeGet(&req)
	} else {
		res, err = messages.ProcessTaskWithResponse[dockerapi.DockerComposeGet, dockerapi.ComposeItem](uint(nodeId), req, defaultTimeout)
	}

	if err != nil {
		return unprocessableEntity(c, err)
	}

	return ok(c, newNodeComposeProjectItemHead(ncp, res))
}

func (h *Handler) GetNodeComposeContainerList(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, errors.New("id should be an integer"))
	}

	ncp, err := h.nodeComposeProjectStore.GetById(uint(nodeId), uint(id))
	if err != nil {
		return unprocessableEntity(c, errors.New("Project not found"))
	}

	req := dockerapi.DockerComposeContainerList{ProjectName: ncp.ProjectName}

	var res *dockerapi.DockerComposeContainerListResponse
	if nodeId == 1 {
		res, err = dockerapi.ComposeContainerList(&req)
	} else {
		res, err = messages.ProcessTaskWithResponse[dockerapi.DockerComposeContainerList, dockerapi.DockerComposeContainerListResponse](uint(nodeId), req, defaultTimeout)
	}

	if err != nil {
		return unprocessableEntity(c, err)
	}
	
	return ok(c, res)
}

func (h *Handler) CreateNodeComposeProject(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	m := model.NodeComposeProject{NodeId: uint(nodeId)}
	r := &nodeComposeProjectCreateRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUnique, err := h.nodeComposeProjectStore.IsUniqueName(uint(nodeId), r.ProjectName)
	if err != nil {
		panic(err)
	}

	if !isUnique {
		return unprocessableEntity(c, duplicateNameError())
	}

	if err := h.nodeComposeProjectStore.Create(&m); err != nil {
		panic(err)
	}

	return created(c, m.Id)
}

func (h *Handler) IsUniqueNodeComposeProjectName(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	value := c.QueryParam("value")

	unique, err := h.nodeComposeProjectStore.IsUniqueName(uint(nodeId), value)
	if err != nil {
		panic(err)
	}

	return ok(c, newUniqueResponse(unique))
}

func (h *Handler) GetNodeComposePull(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, errors.New("id should be an integer"))
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error while upgrading from http to websocket")
		return err
	}
	defer ws.Close()

	ncp, err := h.nodeComposeProjectStore.GetById(uint(nodeId), uint(id))
	if err != nil {
		return unprocessableEntity(c, errors.New("Project not found"))
	}

	clp, err := h.composeLibraryStore.GetByName(ncp.LibraryProjectName)
	if err != nil {
		return unprocessableEntity(c, errors.New("Library Project not found"))
	}

	environmentId := ncp.EnvironmentId
	if ncp.EnvironmentId == nil {
		node, err := h.nodeStore.GetById(uint(nodeId))
		if err != nil {
			return unprocessableEntity(c, errors.New("Node not found"))
		}

		environmentId = node.EnvironmentId
	}

	variables := make(map[string]store.VariableValue)
	if environmentId != nil {
		variables, err = h.variableValueStore.GetMapByEnvironment(*environmentId)
		if err != nil {
			panic(err)			
		}
	}

	req := dockerapi.DockerComposePull{ProjectName: ncp.ProjectName, Definition: clp.Definition, Variables: variables}
	if nodeId == 1 {
		err := dockerapi.ComposePull(&req, ws)
		if err != nil {
			log.Debug().Err(err).Msg("Error while calling ComposePull")
		}
	} else {
		err = messages.ProcessStreamTask[dockerapi.DockerComposePull](uint(nodeId), req, ws)
		if err != nil {
			log.Debug().Err(err).Msg("Error while calling ComposePull ProcessStreamTask")
		}
	}

	return nil
}

func (h *Handler) GetNodeComposeUp(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, errors.New("id should be an integer"))
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error while upgrading from http to websocket")
		return err
	}
	defer ws.Close()

	ncp, err := h.nodeComposeProjectStore.GetById(uint(nodeId), uint(id))
	if err != nil {
		return unprocessableEntity(c, errors.New("Project not found"))
	}

	clp, err := h.composeLibraryStore.GetByName(ncp.LibraryProjectName)
	if err != nil {
		return unprocessableEntity(c, errors.New("Library Project not found"))
	}

	environmentId := ncp.EnvironmentId
	if ncp.EnvironmentId == nil {
		node, err := h.nodeStore.GetById(uint(nodeId))
		if err != nil {
			return unprocessableEntity(c, errors.New("Node not found"))
		}

		environmentId = node.EnvironmentId
	}

	variables := make(map[string]store.VariableValue)
	if environmentId != nil {
		variables, err = h.variableValueStore.GetMapByEnvironment(*environmentId)
		if err != nil {
			panic(err)			
		}
	}

	req := dockerapi.DockerComposeUp{ProjectName: ncp.ProjectName, Definition: clp.Definition, Variables: variables}
	if nodeId == 1 {
		err := dockerapi.ComposeUp(&req, ws)
		if err != nil {
			log.Debug().Err(err).Msg("Error while calling DockerComposeUp")
		}
	} else {
		err = messages.ProcessStreamTask[dockerapi.DockerComposeUp](uint(nodeId), req, ws)
		if err != nil {
			log.Debug().Err(err).Msg("Error while calling DockerComposeUp ProcessStreamTask")
		}
	}

	return nil
}

func (h *Handler) GetNodeComposeDown(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, errors.New("id should be an integer"))
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error while upgrading from http to websocket")
		return err
	}
	defer ws.Close()

	ncp, err := h.nodeComposeProjectStore.GetById(uint(nodeId), uint(id))
	if err != nil {
		return unprocessableEntity(c, errors.New("Project not found"))
	}

	req := dockerapi.DockerComposeDown{ProjectName: ncp.ProjectName}
	if nodeId == 1 {
		err := dockerapi.ComposeDown(&req, ws)
		if err != nil {
			log.Debug().Err(err).Msg("Error while calling DockerComposeDown")
		}
	} else {
		err = messages.ProcessStreamTask[dockerapi.DockerComposeDown](uint(nodeId), req, ws)
		if err != nil {
			log.Debug().Err(err).Msg("Error while calling DockerComposeDown ProcessStreamTask")
		}
	}

	return nil
}

func (h *Handler) GetNodeComposeLogs(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, errors.New("id should be an integer"))
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error while upgrading from http to websocket")
		return err
	}
	defer ws.Close()

	ncp, err := h.nodeComposeProjectStore.GetById(uint(nodeId), uint(id))
	if err != nil {
		return unprocessableEntity(c, errors.New("Project not found"))
	}

	req := dockerapi.DockerComposeLogs{ProjectName: ncp.ProjectName}
	if nodeId == 1 {
		err := dockerapi.ComposeLogs(&req, ws)
		if err != nil {
			log.Debug().Err(err).Msg("Error while calling ComposeLogs")
		}
	} else {
		err = messages.ProcessStreamTask[dockerapi.DockerComposeLogs](uint(nodeId), req, ws)
		if err != nil {
			log.Debug().Err(err).Msg("Error while calling ComposeLogs ProcessStreamTask")
		}
	}

	return nil
}

func (h *Handler) DeleteNodeComposeProject(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, errors.New("id should be an integer"))
	}

	ncp, err := h.nodeComposeProjectStore.GetById(uint(nodeId), uint(id))
	if err != nil {
		return unprocessableEntity(c, errors.New("Project not found"))
	}

	req := dockerapi.DockerComposeDownNoStreaming{ProjectName: ncp.ProjectName}

	if nodeId == 1 {
		err = dockerapi.ComposeDownNoStreaming(&req)
	} else {
		err = messages.ProcessTask[dockerapi.DockerComposeDownNoStreaming](uint(nodeId), req, defaultTimeout)
	}

	if err != nil {
		return unprocessableEntity(c, err)
	}

	err = h.nodeComposeProjectStore.DeleteById(uint(nodeId), uint(id))
	if err != nil {
		return unprocessableEntity(c, err)
	}

	return noContent(c)
}