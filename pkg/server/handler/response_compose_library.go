package handler

import "github.com/productiveops/dokemon/pkg/server/model"

type composeLibraryItemHead struct {
	ProjectName string `json:"projectName"`
}

type composeLibraryItem struct {
	ProjectName string `json:"projectName"`
	Definition  string `json:"definition"`
}

func newComposeLibraryItemHead(m *model.ComposeLibraryItemHead) composeLibraryItemHead {
	return composeLibraryItemHead{ProjectName: m.ProjectName}
}

func newComposeLibraryItemHeadList(rows []model.ComposeLibraryItemHead) []composeLibraryItemHead {
	headRows := make([]composeLibraryItemHead, len(rows))
	for i, r := range rows {
		headRows[i] = newComposeLibraryItemHead(&r)
	}
	return headRows
}

func newComposeLibraryItem(m *model.ComposeLibraryItem) composeLibraryItem {
	return composeLibraryItem{ProjectName: m.ProjectName, Definition: m.Definition}
}
