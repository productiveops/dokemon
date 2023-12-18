package handler

import (
	"github.com/productiveops/dokemon/pkg/server/model"
)

type environmentResponse struct {
	Id   uint		`json:"id"`
	Name string 	`json:"name"`
}

func newEnvironmentResponse(m *model.Environment) *environmentResponse {
	return &environmentResponse{Id: m.Id, Name: m.Name}
}

type environmentHead struct {
	Id   				uint	`json:"id"`
	Name 				string 	`json:"name"`
}

func newEnvironmentHead(m *model.Environment) environmentHead {
	return environmentHead{Id: m.Id, Name: m.Name}
}

func newEnvironmentHeadList(rows []model.Environment) []environmentHead {
	headRows := make([]environmentHead, len(rows))
	for i, r := range rows {
		headRows[i] = newEnvironmentHead(&r)
	}
	return headRows
}
