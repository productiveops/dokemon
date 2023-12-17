package handler

import (
	"dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)

type composeProjectCreateRequest struct {
	ProjectName string `json:"projectName" validate:"required,max=100"`
	Definition  string `json:"definition"`
}

func (r *composeProjectCreateRequest) bind(c echo.Context, m *model.ComposeLibraryItem) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.ProjectName = r.ProjectName
	m.Definition = r.Definition

	return nil
}

type composeProjectUpdateRequest struct {
	ProjectName string `json:"projectName" validate:"required,max=100"`
	NewProjectName string `json:"newProjectName" validate:"required,max=100"`
	Definition  string `json:"definition"`
}

func (r *composeProjectUpdateRequest) bind(c echo.Context, m *model.ComposeLibraryItemUpdate) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.ProjectName = r.ProjectName
	m.NewProjectName = r.NewProjectName
	m.Definition = r.Definition

	return nil
}