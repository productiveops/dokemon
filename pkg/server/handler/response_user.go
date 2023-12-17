package handler

import "dokemon/pkg/server/model"

type userResponse struct {
	Id           uint   `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}

func newUserResponse(m *model.User) *userResponse {
	return &userResponse{
		Id:       m.Id,
		Username: m.UserName,
	}
}

type userHead struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
}

func newUserHeadList(rows []model.User) []userHead {
	headRows := make([]userHead, len(rows))
	for i, r := range rows {
		headRows[i] = userHead{
			Id:       r.Id,
			Username: r.UserName}
	}
	return headRows
}

type userCountResponse struct {
	Count int64   `json:"count"`
}

func newUserCountResponse(count int64) *userCountResponse {
	return &userCountResponse{Count: count}
}
