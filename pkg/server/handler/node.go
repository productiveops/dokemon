package handler

import (
	"encoding/json"
	"errors"
	"slices"
	"strconv"
	"time"

	"github.com/productiveops/dokemon/pkg/crypto"
	"github.com/productiveops/dokemon/pkg/crypto/ske"
	"github.com/productiveops/dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateNode(c echo.Context) error {
	m := model.Node{}
	r := &nodeCreateRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUnique, err := h.nodeStore.IsUniqueName(r.Name)
	if err != nil {
		panic(err)
	}

	if !isUnique {
		return unprocessableEntity(c, duplicateNameError())
	}

	if err := h.nodeStore.Create(&m); err != nil {
		panic(err)
	}

	return created(c, m.Id)
}

func (h *Handler) UpdateNode(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}
	
	m, err := h.nodeStore.GetById(uint(id))
	if err != nil {
		panic(err)
	}
	
	if m == nil {
		return resourceNotFound(c, "Node")
	}	

	r := &nodeUpdateRequest{Id: uint(id)}
	if err := r.bind(c, m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUnique, err := h.nodeStore.IsUniqueNameExcludeItself(r.Name, r.Id)
	if err != nil {
		panic(err)
	}

	if !isUnique {
		return unprocessableEntity(c, duplicateNameError())
	}

	if err := h.nodeStore.Update(m); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) UpdateNodeContainerBaseUrl(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}
	
	m, err := h.nodeStore.GetById(uint(id))
	if err != nil {
		panic(err)
	}
	
	if m == nil {
		return resourceNotFound(c, "Node")
	}

	r := &nodeContainerBaseUrlUpdateRequest{Id: uint(id)}
	if err := r.bind(c, m); err != nil {
		return unprocessableEntity(c, err)
	}
	
	if err := h.nodeStore.UpdateContainerBaseUrl(uint(id), &r.ContainerBaseUrl); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) DeleteNodeById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	if id == 1 {
		return unprocessableEntity(c, errors.New("Node cannot be deleted"))
	}

	exists, err := h.nodeStore.Exists(uint(id))
	if err != nil {
		panic(err)
	}

	if !exists {
		return resourceNotFound(c, "Node")
	}

	if err := h.nodeStore.DeleteById(uint(id)); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) GetNodeList(c echo.Context) error {
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

	rows, totalRows, err := h.nodeStore.GetList(uint(p), uint(s))
	if err != nil {
		panic(err)
	}

	idx := slices.IndexFunc(rows, func(item model.Node) bool { return item.Id == 1 })
	if idx != -1 {
		tokenHash := "-"
		lastPing := time.Now()
		rows[idx].TokenHash = &tokenHash
		rows[idx].LastPing = &lastPing
	}

	return ok(c, newPageResponse[nodeHead](newNodeHeadList(rows), uint(p), uint(s), uint(totalRows)))
}

func (h *Handler) GetNodeHeadById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	m, err := h.nodeStore.GetById(uint(id))
	if err != nil {
		panic(err)
	}

	if m == nil {
		return resourceNotFound(c, "Node")
	}

	if id == 1 {
		tokenHash := "-"
		lastPing := time.Now()
		m.TokenHash = &tokenHash
		m.LastPing = &lastPing
	}

	return ok(c, newNodeHead(m))
}

func (h *Handler) GetNodeById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	m, err := h.nodeStore.GetById(uint(id))
	if err != nil {
		panic(err)
	}

	if m == nil {
		return resourceNotFound(c, "Node")
	}

	if id == 1 {
		tokenHash := "-"
		lastPing := time.Now()
		m.TokenHash = &tokenHash
		m.LastPing = &lastPing
	}

	return ok(c, newNodeResponse(m))
}

func (h *Handler) IsUniqueNodeName(c echo.Context) error {
	value := c.QueryParam("value")

	unique, err := h.nodeStore.IsUniqueName(value)
	if err != nil {
		panic(err)
	}

	return ok(c, newUniqueResponse(unique))
}

func (h *Handler) IsUniqueNodeNameExcludeItself(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	value := c.QueryParam("value")

	unique, err := h.nodeStore.IsUniqueNameExcludeItself(value, uint(id))
	if err != nil {
		panic(err)
	}

	return ok(c, newUniqueResponse(unique))
}

func (h *Handler) GenerateRegistrationToken(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	if id == 1 {
		return unprocessableEntity(c, errors.New("Token cannot be generated for this node"))
	}
	
	m, err := h.nodeStore.GetById(uint(id))
	if err != nil {
		panic(err)
	}
	
	if m == nil {
		return resourceNotFound(c, "Node")
	}

	tokenBytes, _ := json.Marshal(Token{NodeId: uint(id)})
	encryptedToken, _ := ske.Encrypt(string(tokenBytes))

	tokenHash := crypto.HashString(encryptedToken)
	m.TokenHash = &tokenHash
	err = h.nodeStore.Update(m)
	if err != nil {
		panic(err)
	}

	return ok(c, newAgentRegistrationTokenResponse(encryptedToken))
}