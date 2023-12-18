package handler

import (
	"github.com/productiveops/dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)

type composeProjectCreateRequest struct {
	ProjectName string `json:"projectName" validate:"required,max=100"`
	Definition  string `json:"definition"`
}

func (r *composeProjectCreateRequest) bind(c echo.Context, m *model.FileSystemComposeLibraryItem) error {
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

func (r *composeProjectUpdateRequest) bind(c echo.Context, m *model.FileSystemComposeLibraryItemUpdate) error {
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