package handler

import (
	"github.com/productiveops/dokemon/pkg/crypto/ske"
	"github.com/productiveops/dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)

type credentialCreateRequest struct {
	Name      	string  `json:"name" validate:"required,max=50"`
	Service   	*string	`json:"service" validate:"omitempty,max=50"`
	Type      	string  `json:"type" validate:"required,max=50"`
	UserName 	*string `json:"userName" validate:"omitempty,max=100"`
	Secret     	string  `json:"secret" validate:"required"`
}

func (r *credentialCreateRequest) bind(c echo.Context, m *model.Credential) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.Name = r.Name
	m.Service = r.Service
	m.Type = r.Type
	m.UserName = r.UserName
	encryptedSecret, err := ske.Encrypt(r.Secret)
	if err != nil {
		return err
	}
	m.Secret = encryptedSecret

	return nil
}

type credentialUpdateDetailsRequest struct {
	Id        	uint     `json:"id" validate:"required"`
	Name      	string   `json:"name" validate:"required,max=50"`
	Service   	*string	`json:"service" validate:"omitempty,max=50"`
	Type      	string  `json:"type" validate:"required,max=50"`
	UserName 	*string `json:"userName" validate:"omitempty,max=100"` 
}

func (r *credentialUpdateDetailsRequest) bind(c echo.Context, m *model.Credential) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.Id = r.Id
	m.Name = r.Name
	m.Service = r.Service
	m.Type = r.Type
	m.UserName = r.UserName
	
	return nil
}

type credentialUpdateSecretRequest struct {
	Id      uint	`json:"id" validate:"required"`
	Secret	string  `json:"secret" validate:"required"`
}

func (r *credentialUpdateSecretRequest) bind(c echo.Context, m *model.Credential) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	encryptedSecret, err := ske.Encrypt(r.Secret)
	if err != nil {
		return err
	}
	m.Secret = encryptedSecret

	return nil
}