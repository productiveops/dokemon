package handler

import (
	"github.com/productiveops/dokemon/pkg/crypto/ske"
	"github.com/productiveops/dokemon/pkg/server/model"
)

type nodeComposeProjectVariableResponse struct {
	Id       				uint   	`json:"id"`
	Name     				string 	`json:"name"`
	IsSecret 				bool   	`json:"isSecret"`
	Value      				string 	`json:"value"`
}

func newNodeComposeProjectVariableResponse(m *model.NodeComposeProjectVariable) nodeComposeProjectVariableResponse {
	decryptedValue, err := ske.Decrypt(m.Value)
	if err != nil {
		panic(err)
	}

	return nodeComposeProjectVariableResponse{
		Id: m.Id,
		Name: m.Name,
		IsSecret: m.IsSecret,
		Value: decryptedValue,
	}
}

func newNodeComposeProjectVariableList(rows []model.NodeComposeProjectVariable) []nodeComposeProjectVariableResponse {
	headRows := make([]nodeComposeProjectVariableResponse, len(rows))
	for i, r := range rows {
		headRows[i] = newNodeComposeProjectVariableResponse(&r)
	}
	return headRows
}