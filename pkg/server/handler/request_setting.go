package handler

import (
	"dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)


type settingUpdateRequest struct {
	Id        string   `json:"id" validate:"required,max=100"`
	Value     string   `json:"value"`
}

func (r *settingUpdateRequest) bind(c echo.Context, m *model.Setting) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.Id = r.Id
	m.Value = r.Value

	return nil
}