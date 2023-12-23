package handler

import (
	"github.com/productiveops/dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)

type nodeComposeProjectVariableCreateRequest struct {
	NodeComposeProjectId 	uint	`json:"nodeComposeProjectId"`
	Name      				string 	`json:"name" validate:"required,max=100"`
	IsSecret  				bool   	`json:"isSecret"`
	Value      				string 	`json:"value"`
}

func (r *nodeComposeProjectVariableCreateRequest) bind(c echo.Context, m *model.NodeComposeProjectVariable) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.NodeComposeProjectId = r.NodeComposeProjectId
	m.Name = r.Name
	m.IsSecret = r.IsSecret
	m.Value = r.Value

	return nil
}

type nodeComposeProjectVariableUpdateRequest struct {
	Id        				uint   	`json:"id" validate:"required"`
	NodeComposeProjectId 	uint	`json:"nodeComposeProjectId"`
	Name      				string 	`json:"name" validate:"required,max=100"`
	IsSecret  				bool   	`json:"isSecret"`
	Value      				string 	`json:"value"`
}

func (r *nodeComposeProjectVariableUpdateRequest) bind(c echo.Context, m *model.NodeComposeProjectVariable) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.Id = r.Id
	m.NodeComposeProjectId = r.NodeComposeProjectId
	m.Name = r.Name
	m.IsSecret = r.IsSecret
	m.Value = r.Value

	return nil
}
