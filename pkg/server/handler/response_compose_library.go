package handler

import "github.com/productiveops/dokemon/pkg/server/model"

type composeLibraryItemHead struct {
	Id			*uint	`json:"id"`
	ProjectName string 	`json:"projectName"`
	Type 		string	`json:"type"` // github, local
}

type composeLibraryItem struct {
	ProjectName string `json:"projectName"`
	Definition  string `json:"definition"`
}

func newComposeLibraryItemHead(m *model.ComposeLibraryItem) composeLibraryItemHead {
	return composeLibraryItemHead{
		Id: &m.Id,
		ProjectName: m.ProjectName,
		Type: m.Type,
	}
}

func newComposeLibraryItemHeadList(rows []model.ComposeLibraryItem) []composeLibraryItemHead {
	headRows := make([]composeLibraryItemHead, len(rows))
	for i, r := range rows {
		headRows[i] = newComposeLibraryItemHead(&r)
	}
	return headRows
}

func newComposeLibraryItem(m *model.FileSystemComposeLibraryItem) composeLibraryItem {
	return composeLibraryItem{ProjectName: m.ProjectName, Definition: m.Definition}
}
