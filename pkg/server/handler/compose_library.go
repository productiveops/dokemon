package handler

import (
	"sort"
	"strconv"

	"github.com/productiveops/dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)

// Common

func (h *Handler) GetComposeProjectList(c echo.Context) error {
	p, err := strconv.Atoi(c.QueryParam("p"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("p"))
	}

	if p < 1 {
		return unprocessableEntity(c, queryGte1ExpectedError("p"))
	}

	s, err := strconv.Atoi(c.QueryParam("s"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("s"))
	}

	if s < 1 {
		return unprocessableEntity(c, queryGte1ExpectedError("s"))
	}

	rows_local, totalRows, err := h.fileSystemComposeLibraryStore.GetList()
	if err != nil {
		panic(err)
	}

	rows_db, totalRows, err := h.composeLibraryStore.GetList()
	if err != nil {
		panic(err)
	}

	rows := make([]model.ComposeLibraryItem, len(rows_db) + len(rows_local))
	i := 0
	for _, row := range rows_db {
		rows[i] = row
		i++
	}
	for _, row := range rows_local {
		rows[i] = model.ComposeLibraryItem{
			Id: 0,
			CredentialId: nil,
			ProjectName: row.ProjectName,
			Type: "filesystem",
			Url: "",
		}
		i++
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].ProjectName < rows[j].ProjectName
	  })

	return ok(c, newPageResponse[composeLibraryItemHead](newComposeLibraryItemHeadList(rows), uint(p), uint(s), uint(totalRows)))
}

func (h *Handler) isUniqueComposeProjectNameAcrossAllTypes(value string) (bool, error) {
	unique_github, err := h.composeLibraryStore.IsUniqueName(value)
	if err != nil {
		return false, err
	}

	unique_local, err := h.fileSystemComposeLibraryStore.IsUniqueName(value)
	if err != nil {
		return false, err
	}

	return unique_github && unique_local, nil
}

func (h *Handler) isUniqueComposeProjectNameExcludeItselfAcrossAllTypes(newValue string, currentValue string, id uint) (bool, error) {
	var err error

	unique_github := true
	if id != 0 {
		unique_github, err = h.composeLibraryStore.IsUniqueNameExcludeItself(newValue, id)
		if err != nil {
			return false, err
		}
	} else {
		unique_github, err = h.composeLibraryStore.IsUniqueName(newValue)
		if err != nil {
			return false, err
		}
	}

	unique_local := true
	if currentValue != "" {
		unique_local, err = h.fileSystemComposeLibraryStore.IsUniqueNameExcludeItself(newValue, currentValue)
		if err != nil {
			return false, err
		}	
	} else {
		unique_local, err = h.fileSystemComposeLibraryStore.IsUniqueName(newValue)
		if err != nil {
			return false, err
		}	
	}

	return unique_github && unique_local, nil
}

func (h *Handler) IsUniqueComposeProjectName(c echo.Context) error {
	value := c.QueryParam("value")

	unique, err := h.isUniqueComposeProjectNameAcrossAllTypes(value)
	if err != nil {
		panic(err)
	}

	return ok(c, newUniqueResponse(unique))
}

func (h *Handler) IsUniqueComposeProjectNameExcludeItself(c echo.Context) error {
	var err error

	newValue := c.QueryParam("newvalue")
	currentValue := c.QueryParam("currentvalue")
	idstring := c.QueryParam("id")

	var id uint = 0
	if idstring != "" {
		idint, err := strconv.Atoi(idstring)
		if err != nil {
			return unprocessableEntity(c, routeIntExpectedError("id"))
		}
		id = uint(idint)
	}

	unique, err := h.isUniqueComposeProjectNameExcludeItselfAcrossAllTypes(newValue, currentValue, id)
	if err != nil {
		panic(err)
	}

	return ok(c, newUniqueResponse(unique))
}

// File System

func (h *Handler) CreateFileSystemComposeProject(c echo.Context) error {
	m := model.FileSystemComposeLibraryItem{}
	r := &fileSystemComposeProjectCreateRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUnique, err := h.isUniqueComposeProjectNameAcrossAllTypes(m.ProjectName)
	if err != nil {
		panic(err)
	}

	if !isUnique {
		return unprocessableEntity(c, duplicateNameError())
	}

	if err := h.fileSystemComposeLibraryStore.Create(&m); err != nil {
		return unprocessableEntity(c, err)
	}

	return created(c, m.ProjectName)
}

func (h *Handler) UpdateFileSystemComposeProject(c echo.Context) error {
	projectName := c.Param("projectName")

	m := model.FileSystemComposeLibraryItemUpdate{}
	r := &fileSystemComposeProjectUpdateRequest{ProjectName: projectName}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUnique, err := h.isUniqueComposeProjectNameExcludeItselfAcrossAllTypes(m.NewProjectName, m.ProjectName, 0)
	if err != nil {
		panic(err)
	}

	if !isUnique {
		return unprocessableEntity(c, duplicateNameError())
	}

	if err := h.fileSystemComposeLibraryStore.Update(&m); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) DeleteFileSystemComposeProject(c echo.Context) error {
	projectName := c.Param("projectName")

	if err := h.fileSystemComposeLibraryStore.DeleteByName(projectName); err != nil {
		return unprocessableEntity(c, err)
	}

	return noContent(c)
}

func (h *Handler) GetFileSystemComposeProject(c echo.Context) error {
	projectName := c.Param("projectName")

	m, err := h.fileSystemComposeLibraryStore.GetByName(projectName)
	if err != nil {
		return notFound(c, "ComposeProject")
	}

	return ok(c, newFileSystemComposeLibraryItem(m))
}

// GitHub

func (h *Handler) CreateGitHubComposeProject(c echo.Context) error {
	m := model.ComposeLibraryItem{}
	r := &githubComposeProjectCreateRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUnique, err := h.isUniqueComposeProjectNameAcrossAllTypes(m.ProjectName)
	if err != nil {
		panic(err)
	}

	if !isUnique {
		return unprocessableEntity(c, duplicateNameError())
	}

	if err := h.composeLibraryStore.Create(&m); err != nil {
		panic(err)
	}

	return created(c, m.Id)
}

func (h *Handler) UpdateGitHubComposeProject(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}
	
	m, err := h.composeLibraryStore.GetById(uint(id))
	if err != nil {
		panic(err)
	}
	
	if m == nil {
		return resourceNotFound(c, "ComposeLibraryItem")
	}	

	r := &githubComposeProjectUpdateRequest{Id: uint(id)}
	if err := r.bind(c, m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUnique, err := h.isUniqueComposeProjectNameExcludeItselfAcrossAllTypes(m.ProjectName, "", m.Id)
	if err != nil {
		panic(err)
	}

	if !isUnique {
		return unprocessableEntity(c, duplicateNameError())
	}

	if err := h.composeLibraryStore.Update(m); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) DeleteGitHubComposeProject(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	exists, err := h.composeLibraryStore.Exists(uint(id))
	if err != nil {
		panic(err)
	}

	if !exists {
		return resourceNotFound(c, "ComposeLibraryItem")
	}

	if err := h.composeLibraryStore.DeleteById(uint(id)); err != nil {
		return unprocessableEntity(c, err)
	}

	return noContent(c)
}

func (h *Handler) GetGitHubComposeProjectById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	m, err := h.composeLibraryStore.GetById(uint(id))
	if err != nil {
		panic(err)
	}

	if m == nil {
		return resourceNotFound(c, "ComposeLibraryItem")
	}

	return ok(c, newGitHubComposeLibraryItem(m))
}
