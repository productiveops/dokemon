package handler

import (
	"github.com/productiveops/dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)

type nodeCreateRequest struct {
	Name      		string	`json:"name" validate:"required,max=50"`
	EnvironmentId 	*uint	`json:"environmentId"`
}

func (r *nodeCreateRequest) bind(c echo.Context, m *model.Node) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.Name = r.Name
	m.EnvironmentId = r.EnvironmentId

	return nil
}

type nodeUpdateRequest struct {
	Id        		uint    `json:"id" validate:"required"`
	Name      		string  `json:"name" validate:"required,max=50"`
	EnvironmentId 	*uint	`json:"environmentId"`
}

func (r *nodeUpdateRequest) bind(c echo.Context, m *model.Node) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.Id = r.Id
	m.Name = r.Name
	m.EnvironmentId = r.EnvironmentId

	return nil
}

type nodeContainerBaseUrlUpdateRequest struct {
	Id       			uint     `json:"id" validate:"required"`
	ContainerBaseUrl    string   `json:"containerBaseUrl" validate:"max=255"`
}

func (r *nodeContainerBaseUrlUpdateRequest) bind(c echo.Context, m *model.Node) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.Id = r.Id
	m.ContainerBaseUrl = &r.ContainerBaseUrl

	return nil
}