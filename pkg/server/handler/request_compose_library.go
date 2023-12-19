package handler

import (
	"github.com/productiveops/dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)

// File System

type fileSystemComposeProjectCreateRequest struct {
	ProjectName string `json:"projectName" validate:"required,max=100"`
	Definition  string `json:"definition"`
}

func (r *fileSystemComposeProjectCreateRequest) bind(c echo.Context, m *model.FileSystemComposeLibraryItem) error {
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

type fileSystemComposeProjectUpdateRequest struct {
	ProjectName string `json:"projectName" validate:"required,max=100"`
	NewProjectName string `json:"newProjectName" validate:"required,max=100"`
	Definition  string `json:"definition"`
}

func (r *fileSystemComposeProjectUpdateRequest) bind(c echo.Context, m *model.FileSystemComposeLibraryItemUpdate) error {
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

// GitHub

type githubComposeProjectCreateRequest struct {
	ProjectName 	string 	`json:"projectName" validate:"required,max=100"`
	CredentialId  	*uint 	`json:"credentialId"`
	Url 			string 	`json:"url" validate:"required,max=255"`
}

func (r *githubComposeProjectCreateRequest) bind(c echo.Context, m *model.ComposeLibraryItem) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.ProjectName = r.ProjectName
	m.CredentialId = r.CredentialId
	m.Url = r.Url
	m.Type = "github"

	return nil
}

type githubComposeProjectUpdateRequest struct {
	Id        		uint     `json:"id" validate:"required"`
	ProjectName 	string 	`json:"projectName" validate:"required,max=100"`
	CredentialId  	*uint 	`json:"credentialId"`
	Url 			string 	`json:"url" validate:"required,max=255"`
}

func (r *githubComposeProjectUpdateRequest) bind(c echo.Context, m *model.ComposeLibraryItem) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.Id = r.Id
	m.ProjectName = r.ProjectName
	m.CredentialId = r.CredentialId
	m.Url = r.Url
	m.Type = "github"

	return nil
}