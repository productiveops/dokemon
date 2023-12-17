package handler

import (
	"dokemon/pkg/server/model"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateOrUpdateVariableValue(c echo.Context) error {
	variableId, err := strconv.Atoi(c.Param("variableId"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("variableId"))
	}

	environmentId, err := strconv.Atoi(c.Param("environmentId"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("environmentId"))
	}

	variableExists, err := h.variableStore.Exists(uint(variableId))
	if err != nil {
		panic(err)
	}

	if !variableExists {
		return resourceNotFound(c, "Variable")
	}

	environmentExists, err := h.environmentStore.Exists(uint(environmentId))
	if err != nil {
		panic(err)
	}

	if !environmentExists {
		return resourceNotFound(c, "Environment")
	}

	m := model.VariableValue{}
	r := &variableCreateOrUpdateRequest{VariableId: uint(variableId), EnvironmentId: uint(environmentId)}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}
	if err := h.variableValueStore.CreateOrUpdate(&m); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) GetVariableValue(c echo.Context) error {
	variableId, err := strconv.Atoi(c.Param("variableId"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("variableId"))
	}

	environmentId, err := strconv.Atoi(c.Param("environmentId"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("environmentId"))
	}

	variableExists, err := h.variableStore.Exists(uint(variableId))
	if err != nil {
		panic(err)
	}

	if !variableExists {
		return resourceNotFound(c, "Variable")
	}

	environmentExists, err := h.environmentStore.Exists(uint(environmentId))
	if err != nil {
		panic(err)
	}

	if !environmentExists {
		return resourceNotFound(c, "Environment")
	}

	m, err := h.variableValueStore.Get(uint(variableId), uint(environmentId))
	if err != nil {
		panic(err)
	}

	return ok(c, newVariableValueResponse(m))
}
