package handler

import (
	"slices"

	"github.com/productiveops/dokemon/pkg/dockerapi"
	"github.com/productiveops/dokemon/pkg/server/model"
)

type nodeComposeProjectItemHead struct {
	Id                 uint `json:"id"`
	ProjectName        string `json:"projectName"`
	LibraryProjectName string `json:"libraryProjectName"`
	Status             string `json:"status"`
}

func newNodeComposeProjectItemHead(ncp *model.NodeComposeProject, dci *dockerapi.ComposeItem) nodeComposeProjectItemHead {
	res := nodeComposeProjectItemHead{
		Id: ncp.Id,
		ProjectName: ncp.ProjectName,
		LibraryProjectName: ncp.LibraryProjectName,
		Status: "",
	}

	if dci != nil {
		res.Status = dci.Status
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
