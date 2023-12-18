package handler

import "github.com/productiveops/dokemon/pkg/server/model"

type composeLibraryItemHead struct {
	ProjectName string `json:"projectName"`
}

type composeLibraryItem struct {
	ProjectName string `json:"projectName"`
	Definition  string `json:"definition"`
}

func newComposeLibraryItemHead(m *model.LocalComposeLibraryItemHead) composeLibraryItemHead {
	return composeLibraryItemHead{ProjectName: m.ProjectName}
}

func newComposeLibraryItemHeadList(rows []model.LocalComposeLibraryItemHead) []composeLibraryItemHead {
	headRows := make([]composeLibraryItemHead, len(rows))
	for i, r := range rows {
		headRows[i] = newComposeLibraryItemHead(&r)
	}
	return headRows
}

func newComposeLibraryItem(m *model.LocalComposeLibraryItem) composeLibraryItem {
	return composeLibraryItem{ProjectName: m.ProjectName, Definition: m.Definition}
}
