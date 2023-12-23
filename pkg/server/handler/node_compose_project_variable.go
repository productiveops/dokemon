package handler

import (
	"strconv"

	"github.com/productiveops/dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateNodeComposeProjectVariable(c echo.Context) error {
	nodeComposeProjectid, err := strconv.Atoi(c.Param("node_compose_project_id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("node_compose_project_id"))
	}

	m := model.NodeComposeProjectVariable{}
	r := &nodeComposeProjectVariableCreateRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUnique, err := h.nodeComposeProjectVariableStore.IsUniqueName(uint(nodeComposeProjectid), r.Name)
	if err != nil {
		panic(err)
	}

	if !isUnique {
		return unprocessableEntity(c, duplicateNameError())
	}

	if err := h.nodeComposeProjectVariableStore.Create(&m); err != nil {
		panic(err)
	}

	return created(c, m.Id)
}

func (h *Handler) UpdateNodeComposeProjectVariable(c echo.Context) error {
	nodeComposeProjectid, err := strconv.Atoi(c.Param("node_compose_project_id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("node_compose_project_id"))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	exists, err := h.nodeComposeProjectVariableStore.Exists(uint(nodeComposeProjectid), uint(id))
	if err != nil {
		panic(err)
	}

	if !exists {
		return resourceNotFound(c, "NodeComposeProjectVariable")
	}

	m := model.NodeComposeProjectVariable{}
	r := &nodeComposeProjectVariableUpdateRequest{Id: uint(id)}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUnique, err := h.nodeComposeProjectVariableStore.IsUniqueNameExcludeItself(uint(nodeComposeProjectid), r.Name, r.Id)
	if err != nil {
		panic(err)
	}

	if !isUnique {
		return unprocessableEntity(c, duplicateNameError())
	}

	if err := h.nodeComposeProjectVariableStore.Update(&m); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) DeleteNodeComposeProjectVariableById(c echo.Context) error {
	nodeComposeProjectid, err := strconv.Atoi(c.Param("node_compose_project_id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("node_compose_project_id"))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	exists, err := h.nodeComposeProjectVariableStore.Exists(uint(nodeComposeProjectid), uint(id))
	if err != nil {
		panic(err)
	}

	if !exists {
		return resourceNotFound(c, "NodeComposeProjectVariable")
	}

	if err := h.nodeComposeProjectVariableStore.DeleteById(uint(nodeComposeProjectid), uint(id)); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) GetNodeComposeProjectVariableList(c echo.Context) error {
	nodeComposeProjectid, err := strconv.Atoi(c.Param("node_compose_project_id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("node_compose_project_id"))
	}

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

	rows, totalRows, err := h.nodeComposeProjectVariableStore.GetList(uint(nodeComposeProjectid), uint(p), uint(s))
	if err != nil {
		panic(err)
	}

	return ok(c, newPageResponse[nodeComposeProjectVariableResponse](newNodeComposeProjectVariableList(rows), uint(p), uint(s), uint(totalRows)))
}

func (h *Handler) GetNodeComposeProjectVariableById(c echo.Context) error {
	nodeComposeProjectid, err := strconv.Atoi(c.Param("node_compose_project_id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("node_compose_project_id"))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	m, err := h.nodeComposeProjectVariableStore.GetById(uint(nodeComposeProjectid), uint(id))
	if err != nil {
		panic(err)
	}

	if m == nil {
		return resourceNotFound(c, "NodeComposeProjectVariable")
	}

	return ok(c, newNodeComposeProjectVariableResponse(m))
}

func (h *Handler) IsUniqueNodeComposeProjectVariableName(c echo.Context) error {
	nodeComposeProjectid, err := strconv.Atoi(c.Param("node_compose_project_id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("node_compose_project_id"))
	}

	value := c.QueryParam("value")

	unique, err := h.nodeComposeProjectVariableStore.IsUniqueName(uint(nodeComposeProjectid), value)
	if err != nil {
		panic(err)
	}

	return ok(c, newUniqueResponse(unique))
}

func (h *Handler) IsUniqueNodeComposeProjectVariableNameExcludeItself(c echo.Context) error {
	nodeComposeProjectid, err := strconv.Atoi(c.Param("node_compose_project_id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("node_compose_project_id"))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	value := c.QueryParam("value")

	unique, err := h.nodeComposeProjectVariableStore.IsUniqueNameExcludeItself(uint(nodeComposeProjectid), value, uint(id))
	if err != nil {
		panic(err)
	}

	return ok(c, newUniqueResponse(unique))
}