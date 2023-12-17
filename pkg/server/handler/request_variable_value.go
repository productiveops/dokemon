package handler

import (
	"dokemon/pkg/crypto/ske"
	"dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)

type variableCreateOrUpdateRequest struct {
	VariableId    uint   `json:"variableId" validate:"required"`
	EnvironmentId uint   `json:"environmentId" validate:"required"`
	Value         string `json:"value"`
}

func (r *variableCreateOrUpdateRequest) bind(c echo.Context, m *model.VariableValue) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	m.VariableId = r.VariableId
	m.EnvironmentId = r.EnvironmentId

	enryptedValue, err := ske.Encrypt(r.Value)

	if err != nil {
		return err
	}
	m.Value = enryptedValue

	return nil
}