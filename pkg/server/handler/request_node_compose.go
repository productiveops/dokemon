package handler

import (
	"github.com/productiveops/dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)

type nodeComposeProjectAddFromLibraryRequest struct {
	ProjectName      	string  `json:"projectName" validate:"required,max=50"`
	LibraryProjectId	*uint	`json:"libraryProjectId"`
	LibraryProjectName	string  `json:"libraryProjectName" validate:"required,max=50"`
}

func (r *nodeComposeProjectAddFromLibraryRequest) bind(c echo.Context, m *model.NodeComposeProject) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.LibraryProjectId = r.LibraryProjectId
	m.ProjectName = r.ProjectName
	m.LibraryProjectName = &r.LibraryProjectName

	return nil
}

type nodeComposeGitHubCreateRequest struct {
	ProjectName     string  `json:"projectName" validate:"required,max=50"`
	CredentialId  	*uint 	`json:"credentialId"`
	Url 			string 	`json:"url" validate:"required,max=255"`
}

func (r *nodeComposeGitHubCreateRequest) bind(c echo.Context, m *model.NodeComposeProject) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.Type = "github"
	m.ProjectName = r.ProjectName
	m.Url = &r.Url
	m.CredentialId = r.CredentialId

	return nil
}

type nodeComposeGitHubUpdateRequest struct {
	Id        		uint    `json:"id" validate:"required"`
	ProjectName     string  `json:"projectName" validate:"required,max=50"`
	CredentialId  	*uint 	`json:"credentialId"`
	Url 			string 	`json:"url" validate:"required,max=255"`
}

func (r *nodeComposeGitHubUpdateRequest) bind(c echo.Context, m *model.NodeComposeProject) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.Id = r.Id
	m.Type = "github"
	m.ProjectName = r.ProjectName
	m.Url = &r.Url
	m.CredentialId = r.CredentialId

	return nil
}

type nodeComposeLocalCreateRequest struct {
	ProjectName     string  `json:"projectName" validate:"required,max=50"`
	Definition  	string 	`json:"definition"`
}

func (r *nodeComposeLocalCreateRequest) bind(c echo.Context, m *model.NodeComposeProject) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.Type = "local"
	m.ProjectName = r.ProjectName
	m.Definition = &r.Definition

	return nil
}

type nodeComposeLocalUpdateRequest struct {
	Id        		uint    `json:"id" validate:"required"`
	ProjectName     string  `json:"projectName" validate:"required,max=50"`
	Definition  	string 	`json:"definition"`
}

func (r *nodeComposeLocalUpdateRequest) bind(c echo.Context, m *model.NodeComposeProject) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.Id = r.Id
	m.Type = "local"
	m.ProjectName = r.ProjectName
	m.Definition = &r.Definition

	return nil
}