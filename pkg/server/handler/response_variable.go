package handler

import "github.com/productiveops/dokemon/pkg/server/model"

type variableResponse struct {
	Id       uint   `json:"id"`
	Name     string `json:"name"`
	IsSecret bool   `json:"isSecret"`
}

type variableHead struct {
	Id       uint              `json:"id"`
	Name     string            `json:"name"`
	IsSecret bool              `json:"isSecret"`
	Values   map[string]string `json:"values"`
}

func newVariableResponse(m *model.Variable) *variableResponse {
	return &variableResponse{
		Id: m.Id,
		Name: m.Name,
		IsSecret: m.IsSecret,
	}
}