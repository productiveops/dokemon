package handler

import "github.com/productiveops/dokemon/pkg/server/model"

type composeLibraryItemHead struct {
	Id			*uint	`json:"id"`
	ProjectName string 	`json:"projectName"`
	Type 		string	`json:"type"` // github, local
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

type fileSystemComposeLibraryItem struct {
	ProjectName string `json:"projectName"`
	Definition  string `json:"definition"`
}

func newFileSystemComposeLibraryItem(m *model.FileSystemComposeLibraryItem) fileSystemComposeLibraryItem {
	return fileSystemComposeLibraryItem{
		ProjectName: m.ProjectName,
		Definition: m.Definition,
	}
}

type gitHubComposeLibraryItem struct {
	Id				uint	`json:"id"`
	CredentialId	*uint	`json:"credentialId"`
	ProjectName 	string 	`json:"projectName"`
	Url 			string 	`json:"url"`
}

func newGitHubComposeLibraryItem(m *model.ComposeLibraryItem) gitHubComposeLibraryItem {
	return gitHubComposeLibraryItem{
		Id: m.Id,
		CredentialId: m.CredentialId,
		ProjectName: m.ProjectName,
		Url: m.Url,
	}
}
