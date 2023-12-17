package handler

import (
	"dokemon/pkg/server/model"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateVariable(c echo.Context) error {
	m := model.Variable{}
	r := &variableCreateRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUnique, err := h.variableStore.IsUniqueName(r.Name)
	if err != nil {
		panic(err)
	}

	if !isUnique {
		return unprocessableEntity(c, duplicateNameError())
	}

	if err := h.variableStore.Create(&m); err != nil {
		panic(err)
	}

	return created(c, m.Id)
}

func (h *Handler) UpdateVariable(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	exists, err := h.variableStore.Exists(uint(id))
	if err != nil {
		panic(err)
	}

	if !exists {
		return resourceNotFound(c, "Variable")
	}

	m := model.Variable{}
	r := &variableUpdateRequest{Id: uint(id)}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUnique, err := h.variableStore.IsUniqueNameExcludeItself(r.Name, r.Id)
	if err != nil {
		panic(err)
	}

	if !isUnique {
		return unprocessableEntity(c, duplicateNameError())
	}

	if err := h.variableStore.Update(&m); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) DeleteVariableById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	exists, err := h.variableStore.Exists(uint(id))
	if err != nil {
		panic(err)
	}

	if !exists {
		return resourceNotFound(c, "Variable")
	}

	if err := h.variableStore.DeleteById(uint(id)); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) GetVariableList(c echo.Context) error {
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

	rows, totalRows, err := h.variableStore.GetList(uint(p), uint(s))
	if err != nil {
		panic(err)
	}

	ret := make([]variableHead, len(rows))
	for i, row := range rows {
		valuesmap, err := h.variableValueStore.GetMap(row.Id)
		if err != nil {
			panic(err)
		}
		ret[i] = variableHead{Id: row.Id, Name: row.Name, IsSecret: row.IsSecret, Values: valuesmap}
	}

	return ok(c, newPageResponse[variableHead](ret, uint(p), uint(s), uint(totalRows)))
}

func (h *Handler) GetVariableById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	m, err := h.variableStore.GetById(uint(id))
	if err != nil {
		panic(err)
	}

	if m == nil {
		return resourceNotFound(c, "Variable")
	}

	return ok(c, newVariableResponse(m))
}

func (h *Handler) IsUniqueVariableName(c echo.Context) error {
	value := c.QueryParam("value")

	unique, err := h.variableStore.IsUniqueName(value)
	if err != nil {
		panic(err)
	}

	return ok(c, newUniqueResponse(unique))
}

func (h *Handler) IsUniqueVariableNameExcludeItself(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	value := c.QueryParam("value")

	unique, err := h.variableStore.IsUniqueNameExcludeItself(value, uint(id))
	if err != nil {
		panic(err)
	}

	return ok(c, newUniqueResponse(unique))
}