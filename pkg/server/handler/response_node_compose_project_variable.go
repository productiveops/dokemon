package handler

import "github.com/productiveops/dokemon/pkg/server/model"

type nodeComposeProjectVariableResponse struct {
	Id       				uint   	`json:"id"`
	NodeComposeProjectId 	uint   	`json:"nodeComposeProjectId"`
	Name     				string 	`json:"name"`
	IsSecret 				bool   	`json:"isSecret"`
	Value      				string 	`json:"value"`
}

func newNodeComposeProjectVariableResponse(m *model.NodeComposeProjectVariable) *nodeComposeProjectVariableResponse {
	return &nodeComposeProjectVariableResponse{
		Id: m.Id,
		NodeComposeProjectId: m.NodeComposeProjectId,
		Name: m.Name,
		IsSecret: m.IsSecret,
		Value: m.Value,
	}
}