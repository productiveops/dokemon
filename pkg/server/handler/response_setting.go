package handler

import (
	"github.com/productiveops/dokemon/pkg/server/model"
)

type settingResponse struct {
	Id   	string	`json:"id"`
	Value 	string 	`json:"value"`
}

func newSettingResponse(m *model.Setting) *settingResponse {
	return &settingResponse{Id: m.Id, Value: m.Value}
}
