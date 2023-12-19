package handler

import (
	"github.com/labstack/echo/v4"
)

type gitHubfileContentRetrieveRequest struct {
	CredentialId	*uint `json:"credentialId"`
	Url     		string	`json:"url" validate:"required"`
}

func (r *gitHubfileContentRetrieveRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	return nil
}
