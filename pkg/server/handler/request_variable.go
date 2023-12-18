package handler

import (
	"github.com/productiveops/dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)

type variableCreateRequest struct {
	Name      string `json:"name" validate:"required,max=100"`
	IsSecret  bool   `json:"isSecret"`
}

func (r *variableCreateRequest) bind(c echo.Context, m *model.Variable) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.Name = r.Name
	m.IsSecret = r.IsSecret

	return nil
}

type variableUpdateRequest struct {
	Id        uint   `json:"id" validate:"required"`
	Name      string `json:"name" validate:"required,max=100"`
	IsSecret  bool   `json:"isSecret"`
}

func (r *variableUpdateRequest) bind(c echo.Context, m *model.Variable) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.Id = r.Id
	m.Name = r.Name
	m.IsSecret = r.IsSecret

	return nil
}
