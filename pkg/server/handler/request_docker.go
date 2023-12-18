package handler

import (
	"github.com/productiveops/dokemon/pkg/dockerapi"

	"github.com/labstack/echo/v4"
)

type dockerContainerStartRequest struct {
	Id string `json:"id" validate:"required,max=100"`
}

func (r *dockerContainerStartRequest) bind(c echo.Context, m *dockerapi.DockerContainerStart) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	
	if err := c.Validate(r); err != nil {
		return err
	}

	m.Id = r.Id
	return nil
}

type dockerContainerStopRequest struct {
	Id string `json:"id" validate:"required,max=100"`
}

func (r *dockerContainerStopRequest) bind(c echo.Context, m *dockerapi.DockerContainerStop) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	
	if err := c.Validate(r); err != nil {
		return err
	}

	m.Id = r.Id
	return nil
}

type dockerContainerRestartRequest struct {
	Id string `json:"id" validate:"required,max=100"`
}

func (r *dockerContainerRestartRequest) bind(c echo.Context, m *dockerapi.DockerContainerRestart) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	
	if err := c.Validate(r); err != nil {
		return err
	}

	m.Id = r.Id
	return nil
}

type dockerContainerRemoveRequest struct {
	Id     	string	`json:"id" validate:"required,max=100"`
	Force   bool	`json:"force" validate:"required"`
}

func (r *dockerContainerRemoveRequest) bind(c echo.Context, m *dockerapi.DockerContainerRemove) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	
	if err := c.Validate(r); err != nil {
		return err
	}

	m.Id = r.Id
	m.Force = r.Force
	return nil
}

type dockerImageRemoveRequest struct {
	Id      string  `json:"id" validate:"required,max=100"`
	Force   bool	`json:"force" validate:"required"`
}

func (r *dockerImageRemoveRequest) bind(c echo.Context, m *dockerapi.DockerImageRemove) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	
	if err := c.Validate(r); err != nil {
		return err
	}

	m.Id = r.Id
	m.Force = r.Force
	return nil
}

type dockerImagesPruneRequest struct {
	All      bool  `json:"all" validate:"required"`
}

func (r *dockerImagesPruneRequest) bind(c echo.Context, m *dockerapi.DockerImagesPrune) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	
	if err := c.Validate(r); err != nil {
		return err
	}

	m.All = r.All
	return nil
}

type dockerVolumeRemoveRequest struct {
	Name      string  `json:"name" validate:"required,max=200"`
}

func (r *dockerVolumeRemoveRequest) bind(c echo.Context, m *dockerapi.DockerVolumeRemove) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	
	if err := c.Validate(r); err != nil {
		return err
	}

	m.Name = r.Name
	return nil
}

type dockerVolumePruneRequest struct {
	All      bool  `json:"all" validate:"required"`
}

func (r *dockerVolumePruneRequest) bind(c echo.Context, m *dockerapi.DockerVolumesPrune) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	
	if err := c.Validate(r); err != nil {
		return err
	}

	m.All = r.All
	return nil
}

type dockerNetworkRemoveRequest struct {
	Id      string  `json:"id" validate:"required,max=100"`
}

func (r *dockerNetworkRemoveRequest) bind(c echo.Context, m *dockerapi.DockerNetworkRemove) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	
	if err := c.Validate(r); err != nil {
		return err
	}

	m.Id = r.Id
	return nil
}

type dockerComposeProjectCreateRequest struct {
	ProjectName string `json:"projectName" validate:"required,max=100"`
	Definition string `json:"definition"`
}

func (r *dockerComposeProjectCreateRequest) bind(c echo.Context, m *dockerapi.DockerComposeProjectCreate) error {
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

type dockerComposeProjectUpdateRequest struct {
	ProjectName string `json:"projectName" validate:"required,max=100"`
	Definition string `json:"definition"`
}

func (r *dockerComposeProjectUpdateRequest) bind(c echo.Context, m *dockerapi.DockerComposeProjectUpdate) error {
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