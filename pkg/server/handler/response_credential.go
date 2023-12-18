package handler

import (
	"github.com/productiveops/dokemon/pkg/server/model"
)

type credentialResponse struct {
	Id   		uint	`json:"id"`
	Name 		string 	`json:"name"`
	Service 	*string `json:"service"`
	Type 		string 	`json:"type"`
	UserName 	*string `json:"userName"`
}

func newCredentialResponse(m *model.Credential) *credentialResponse {
	return &credentialResponse{
		Id: m.Id,
		Name: m.Name,
		Service: m.Service,
		Type: m.Type,
		UserName: m.UserName,
	}
}

type credentialHead struct {
	Id   		uint	`json:"id"`
	Name 		string 	`json:"name"`
	Service 	*string `json:"service"`
	Type 		string 	`json:"type"`
	UserName 	*string	`json:"userName"`
}

func newCredentialHead(m *model.Credential) credentialHead {
	return credentialHead{
		Id: m.Id,
		Name: m.Name,
		Service: m.Service,
		Type: m.Type,
		UserName: m.UserName,
	}
}

func newCredentialHeadList(rows []model.Credential) []credentialHead {
	headRows := make([]credentialHead, len(rows))
	for i, r := range rows {
		headRows[i] = newCredentialHead(&r)
	}
	return headRows
}
