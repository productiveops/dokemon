package handler

import (
	"errors"
	"strconv"

	"github.com/productiveops/dokemon/pkg/crypto/ske"
	"github.com/productiveops/dokemon/pkg/dockerapi"
	"github.com/productiveops/dokemon/pkg/messages"
	"github.com/productiveops/dokemon/pkg/server/model"
	"github.com/productiveops/dokemon/pkg/server/store"

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

	if ncp.LibraryProjectId != nil || ncp.LibraryProjectName != nil {
		definition, credentialId, url, err := h.getComposeProjectDefinitionFromLibrary(ncp)
		if err != nil {
			panic(err)
		}
		ncp.Definition = &definition
		ncp.CredentialId = credentialId
		ncp.Url = url
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

	return ok(c, newNodeComposeProjectItem(ncp, res))
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

func (h *Handler) CreateGitHubNodeComposeProject(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	m := model.NodeComposeProject{NodeId: uint(nodeId)}
	r := &nodeComposeGitHubCreateRequest{}
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

func (h *Handler) CreateLocalNodeComposeProject(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	m := model.NodeComposeProject{NodeId: uint(nodeId)}
	r := &nodeComposeLocalCreateRequest{}
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

func (h *Handler) UpdateGitHubNodeComposeProject(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}
	
	m, err := h.nodeComposeProjectStore.GetById(uint(nodeId), uint(id))
	if err != nil {
		panic(err)
	}
	
	if m == nil {
		return resourceNotFound(c, "NodeComposeProject")
	}	

	r := &nodeComposeGitHubUpdateRequest{Id: uint(id)}
	if err := r.bind(c, m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUnique, err := h.nodeComposeProjectStore.IsUniqueNameExcludeItself(uint(nodeId), r.ProjectName, r.Id)
	if err != nil {
		panic(err)
	}

	if !isUnique {
		return unprocessableEntity(c, duplicateNameError())
	}

	if m.LibraryProjectId != nil || m.LibraryProjectName != nil {
		m.CredentialId = nil
		m.Url = nil
	}

	if err := h.nodeComposeProjectStore.Update(m); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) UpdateLocalNodeComposeProject(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}
	
	m, err := h.nodeComposeProjectStore.GetById(uint(nodeId), uint(id))
	if err != nil {
		panic(err)
	}
	
	if m == nil {
		return resourceNotFound(c, "NodeComposeProject")
	}	

	r := &nodeComposeLocalUpdateRequest{Id: uint(id)}
	if err := r.bind(c, m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUnique, err := h.nodeComposeProjectStore.IsUniqueNameExcludeItself(uint(nodeId), r.ProjectName, r.Id)
	if err != nil {
		panic(err)
	}

	if !isUnique {
		return unprocessableEntity(c, duplicateNameError())
	}

	if m.LibraryProjectId != nil || m.LibraryProjectName != nil {
		m.Definition = nil
	}

	if err := h.nodeComposeProjectStore.Update(m); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) AddNodeComposeProjectFromLibrary(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	m := model.NodeComposeProject{NodeId: uint(nodeId)}
	r := &nodeComposeProjectAddFromLibraryRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	if r.LibraryProjectId != nil {
		m.Type = "github"
		_, err := h.composeLibraryStore.GetById(*r.LibraryProjectId)
		if err != nil {
			return unprocessableEntity(c, errors.New("Library Project not found"))
		}
	} else {
		m.Type = "local"
		_, err := h.fileSystemComposeLibraryStore.GetByName(r.LibraryProjectName)
		if err != nil {
			return unprocessableEntity(c, errors.New("Library Project not found"))
		}
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

func (h *Handler) IsUniqueNodeComposeProjectNameExcludeItself(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.Param("nodeId"))
	if err != nil {
		return unprocessableEntity(c, errors.New("nodeId should be an integer"))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, errors.New("id should be an integer"))
	}

	value := c.QueryParam("value")

	unique, err := h.nodeComposeProjectStore.IsUniqueNameExcludeItself(uint(nodeId), value, uint(id))
	if err != nil {
		panic(err)
	}

	return ok(c, newUniqueResponse(unique))
}

func (h *Handler) getComposeProjectDefinitionFromLibrary(ncp *model.NodeComposeProject) (definition string, credentialId *uint, url *string, err error) {
	if ncp.LibraryProjectId == nil {
		clp, err := h.fileSystemComposeLibraryStore.GetByName(*ncp.LibraryProjectName)
		if err != nil {
			return "", nil, nil, errors.New("Library Project not found")
		}	
		definition = clp.Definition
		credentialId = nil
		url = nil
	} else {
		gclp, err := h.composeLibraryStore.GetById(*ncp.LibraryProjectId)
		if err != nil {
			return "", nil, nil, errors.New("Library Project not found")
		}

		decryptedSecret := ""
		if gclp.CredentialId != nil {
			credential, err := h.credentialStore.GetById(*gclp.CredentialId)
			if err != nil {
				return "", nil, nil, errors.New("Credentials not found")
			}
			
			decryptedSecret, err = ske.Decrypt(credential.Secret)
			if err != nil {
				panic(err)
			}
		}

		content, err := getGitHubFileContent(gclp.Url, decryptedSecret)
		if err != nil {
			return "", nil, nil, errors.New("Error while retrieving file content from GitHub")
		}
		
		definition = content
		credentialId = gclp.CredentialId
		url = &gclp.Url
	}

	return definition, credentialId, url, nil
}

func (h *Handler) getComposeProjectDefinition(ncp *model.NodeComposeProject) (string, error) {
	var err error
	definition := ""

	if ncp.LibraryProjectId != nil || ncp.LibraryProjectName != nil {
		definition, _, _, err = h.getComposeProjectDefinitionFromLibrary(ncp)
		if err != nil {
			return "", err
		}
	} else if ncp.Type == "local" {
		definition = *ncp.Definition
	} else if ncp.Type == "github" {
		decryptedSecret := ""
		if ncp.CredentialId != nil {
			credential, err := h.credentialStore.GetById(*ncp.CredentialId)
			if err != nil {
				return "", errors.New("Credentials not found")
			}
			
			decryptedSecret, err = ske.Decrypt(credential.Secret)
			if err != nil {
				panic(err)
			}
		}

		content, err := getGitHubFileContent(*ncp.Url, decryptedSecret)
		if err != nil {
			return "", errors.New("Error while retrieving file content from GitHub")
		}
		
		definition = content
	}

	return definition, nil
}

func (h *Handler) getComposeVariables(environmentId *uint, nodeComposeProjectId uint) map[string]store.VariableValue {
	var err error

	variables := make(map[string]store.VariableValue)
	if environmentId != nil {
		variables, err = h.variableValueStore.GetMapByEnvironment(*environmentId)
		if err != nil {
			panic(err)
		}
	}

	composeVariables, _, err := h.nodeComposeProjectVariableStore.GetList(nodeComposeProjectId, 1, 10000)
	if err != nil {
		panic(err)
	}

	for _, v := range composeVariables {
		decryptedValue, err := ske.Decrypt(v.Value)
		if err != nil {
			panic(err)
		}

		variables[v.Name] = store.VariableValue{
			IsSecret: v.IsSecret,
			Value: &decryptedValue,
		}
	}

	return variables
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

	ncp, err := h.nodeComposeProjectStore.GetById(uint(nodeId), uint(id))
	if err != nil {
		return unprocessableEntity(c, errors.New("Project not found"))
	}

	definition, err := h.getComposeProjectDefinition(ncp)
	if err != nil {
		return unprocessableEntity(c, err)
	}

	environmentId := ncp.EnvironmentId
	if ncp.EnvironmentId == nil {
		node, err := h.nodeStore.GetById(uint(nodeId))
		if err != nil {
			return unprocessableEntity(c, errors.New("Node not found"))
		}

		environmentId = node.EnvironmentId
	}

	variables := h.getComposeVariables(environmentId, uint(id))

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error while upgrading from http to websocket")
		return err
	}
	defer ws.Close()

	req := dockerapi.DockerComposePull{ProjectName: ncp.ProjectName, Definition: definition, Variables: variables}
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
	
	ncp, err := h.nodeComposeProjectStore.GetById(uint(nodeId), uint(id))
	if err != nil {
		return unprocessableEntity(c, errors.New("Project not found"))
	}

	definition, err := h.getComposeProjectDefinition(ncp)
	if err != nil {
		return unprocessableEntity(c, err)
	}

	environmentId := ncp.EnvironmentId
	if ncp.EnvironmentId == nil {
		node, err := h.nodeStore.GetById(uint(nodeId))
		if err != nil {
			return unprocessableEntity(c, errors.New("Node not found"))
		}

		environmentId = node.EnvironmentId
	}

	variables := h.getComposeVariables(environmentId, uint(id))

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error while upgrading from http to websocket")
		return err
	}
	defer ws.Close()

	req := dockerapi.DockerComposeUp{ProjectName: ncp.ProjectName, Definition: definition, Variables: variables}
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

	ncp, err := h.nodeComposeProjectStore.GetById(uint(nodeId), uint(id))
	if err != nil {
		return unprocessableEntity(c, errors.New("Project not found"))
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error while upgrading from http to websocket")
		return err
	}
	defer ws.Close()

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