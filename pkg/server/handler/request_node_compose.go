package handler

import (
	"github.com/productiveops/dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)

type nodeComposeProjectCreateRequest struct {
	ProjectName      	string  `json:"projectName" validate:"required,max=50"`
	LibraryProjectId	*uint	`json:"libraryProjectId"`
	LibraryProjectName	string  `json:"libraryProjectName" validate:"required,max=50"`
}

func (r *nodeComposeProjectCreateRequest) bind(c echo.Context, m *model.NodeComposeProject) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.LibraryProjectId = r.LibraryProjectId
	m.ProjectName = r.ProjectName
	m.LibraryProjectName = r.LibraryProjectName

	return nil
}