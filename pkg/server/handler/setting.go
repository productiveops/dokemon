package handler

import (
	"github.com/labstack/echo/v4"
)

func (h *Handler) UpdateSetting(c echo.Context) error {
	id := c.Param("id")

	m, err := h.settingStore.GetById(id)
	if err != nil {
		panic(err)
	}
	
	if m == nil {
		return resourceNotFound(c, "Setting")
	}	

	r := &settingUpdateRequest{Id: id}
	if err := r.bind(c, m); err != nil {
		return unprocessableEntity(c, err)
	}

	if err := h.settingStore.Update(m); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) GetSettingById(c echo.Context) error {
	id := c.Param("id")

	m, err := h.settingStore.GetById(id)
	if err != nil {
		panic(err)
	}

	if m == nil {
		return resourceNotFound(c, "Setting")
	}

	return ok(c, newSettingResponse(m))
}
