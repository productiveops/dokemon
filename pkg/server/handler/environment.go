package handler

import (
	"errors"
	"strconv"

	"github.com/productiveops/dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateEnvironment(c echo.Context) error {
	m := model.Environment{}
	r := &environmentCreateRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUnique, err := h.environmentStore.IsUniqueName(r.Name)
	if err != nil {
		panic(err)
	}

	if !isUnique {
		return unprocessableEntity(c, duplicateNameError())
	}

	if err := h.environmentStore.Create(&m); err != nil {
		panic(err)
	}

	return created(c, m.Id)
}

func (h *Handler) UpdateEnvironment(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}
	
	m, err := h.environmentStore.GetById(uint(id))
	if err != nil {
		panic(err)
	}
	
	if m == nil {
		return resourceNotFound(c, "Environment")
	}	

	r := &environmentUpdateRequest{Id: uint(id)}
	if err := r.bind(c, m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUnique, err := h.environmentStore.IsUniqueNameExcludeItself(r.Name, r.Id)
	if err != nil {
		panic(err)
	}

	if !isUnique {
		return unprocessableEntity(c, duplicateNameError())
	}

	if err := h.environmentStore.Update(m); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) DeleteEnvironmentById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	if id == 1 {
		return unprocessableEntity(c, errors.New("Environment cannot be deleted"))
	}

	exists, err := h.environmentStore.Exists(uint(id))
	if err != nil {
		panic(err)
	}

	if !exists {
		return resourceNotFound(c, "Environment")
	}

	if err := h.environmentStore.DeleteById(uint(id)); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) GetEnvironmentList(c echo.Context) error {
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

	rows, totalRows, err := h.environmentStore.GetList(uint(p), uint(s))
	if err != nil {
		panic(err)
	}

	return ok(c, newPageResponse[environmentHead](newEnvironmentHeadList(rows), uint(p), uint(s), uint(totalRows)))
}

func (h *Handler) GetEnvironmentHeadById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	m, err := h.environmentStore.GetById(uint(id))
	if err != nil {
		panic(err)
	}

	if m == nil {
		return resourceNotFound(c, "Environment")
	}

	return ok(c, newEnvironmentHead(m))
}

func (h *Handler) GetEnvironmentById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	m, err := h.environmentStore.GetById(uint(id))
	if err != nil {
		panic(err)
	}

	if m == nil {
		return resourceNotFound(c, "Environment")
	}

	return ok(c, newEnvironmentResponse(m))
}

func (h *Handler) IsUniqueEnvironmentName(c echo.Context) error {
	value := c.QueryParam("value")

	unique, err := h.environmentStore.IsUniqueName(value)
	if err != nil {
		panic(err)
	}

	return ok(c, newUniqueResponse(unique))
}

func (h *Handler) IsUniqueEnvironmentNameExcludeItself(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	value := c.QueryParam("value")

	unique, err := h.environmentStore.IsUniqueNameExcludeItself(value, uint(id))
	if err != nil {
		panic(err)
	}

	return ok(c, newUniqueResponse(unique))
}

func (h *Handler) GetEnvironmentMap(c echo.Context) error {
	m, err := h.environmentStore.GetMap()
	if err != nil {
		panic(err)
	}

	return ok(c, m)
}