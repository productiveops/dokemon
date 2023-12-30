package handler

import (
	"slices"

	"github.com/productiveops/dokemon/pkg/dockerapi"
	"github.com/productiveops/dokemon/pkg/server/model"
)

type nodeComposeProjectItemHead struct {
	Id                 uint `json:"id"`
	ProjectName        string `json:"projectName"`
	Type        	   string `json:"type"`
	LibraryProjectId   *uint `json:"libraryProjectId"`
	LibraryProjectName *string `json:"libraryProjectName"`
	Status             string `json:"status"`
	Stale             	string `json:"stale"`
}

func newNodeComposeProjectItemHead(ncp *model.NodeComposeProject, dci *dockerapi.ComposeItem) nodeComposeProjectItemHead {
	res := nodeComposeProjectItemHead{
		Id: ncp.Id,
		ProjectName: ncp.ProjectName,
		Type: ncp.Type,
		LibraryProjectId: ncp.LibraryProjectId,
		LibraryProjectName: ncp.LibraryProjectName,
		Status: "",
		Stale: "",
	}

	if dci != nil {
		res.Status = dci.Status
		res.Stale = dci.Stale
	}

	return res
}

func newNodeComposeProjectItemList(ncplist []model.NodeComposeProject, dcilist []dockerapi.ComposeItem) []nodeComposeProjectItemHead {
	res := make([]nodeComposeProjectItemHead, len(ncplist))

	for i, ncp := range ncplist {
		idx := slices.IndexFunc(dcilist, func(dci dockerapi.ComposeItem) bool { return dci.Name == ncp.ProjectName } )
		var dci *dockerapi.ComposeItem = nil
		if idx != -1 {
			dci = &dcilist[idx]
		}
		res[i] = newNodeComposeProjectItemHead(&ncp, dci)
	}

	return res
}

type nodeComposeProjectItem struct {
	Id                 	uint `json:"id"`
	ProjectName        	string `json:"projectName"`
	Type        	   	string `json:"type"`
	LibraryProjectId   	*uint `json:"libraryProjectId"`
	LibraryProjectName 	*string `json:"libraryProjectName"`
	Url        	   		*string `json:"url"`
	CredentialId       	*uint `json:"credentialId"`
	Definition       	*string `json:"definition"`
	Status             	string `json:"status"`
	Stale             	string `json:"stale"`
}

func newNodeComposeProjectItem(ncp *model.NodeComposeProject, dci *dockerapi.ComposeItem) nodeComposeProjectItem {
	res := nodeComposeProjectItem{
		Id: ncp.Id,
		ProjectName: ncp.ProjectName,
		Type: ncp.Type,
		LibraryProjectId: ncp.LibraryProjectId,
		LibraryProjectName: ncp.LibraryProjectName,
		Url: ncp.Url,
		CredentialId: ncp.CredentialId,
		Definition: ncp.Definition,
		Status: "",
		Stale: "",
	}

	if dci != nil {
		res.Status = dci.Status
		res.Stale = dci.Stale
	}

	return res
}
