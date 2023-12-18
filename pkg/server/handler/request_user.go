package handler

import (
	"github.com/productiveops/dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)

type userCreateRequest struct {
	UserName    string `json:"userName" validate:"max=255"`
	Password 	string `json:"password" validate:"min=8,max=255"`
}

func (r *userCreateRequest) bind(c echo.Context, m *model.User) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.UserName = r.UserName
	m.PasswordHash = r.Password

	return nil
}

type userLoginRequest struct {
	UserName    string `json:"userName" validate:"max=255"`
	Password 	string `json:"password" validate:"max=255"`
}

func (r *userLoginRequest) bind(c echo.Context, m *model.User) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.UserName = r.UserName
	m.PasswordHash = r.Password

	return nil
}

type userUpdateRequest struct {
	Id           uint   `json:"id" validate:"required"`
	UserName     string `json:"userName" validate:"max=255"`
}

func (r *userUpdateRequest) bind(c echo.Context, m *model.User) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.Id = r.Id
	m.UserName = r.UserName

	return nil
}


type userChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"max=255"`
	NewPassword     string `json:"newPassword" validate:"min=8,max=255"`
}

func (r *userChangePasswordRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	return nil
}