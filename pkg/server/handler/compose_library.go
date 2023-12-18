package handler

import (
	"sort"
	"strconv"

	"github.com/productiveops/dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)

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

func (h *Handler) CreateFileSystemComposeProject(c echo.Context) error {
	m := model.FileSystemComposeLibraryItem{}
	r := &composeProjectCreateRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	if err := h.fileSystemComposeLibraryStore.Create(&m); err != nil {
		panic(err)
	}

	return created(c, m.ProjectName)
}

func (h *Handler) UpdateFileSystemComposeProject(c echo.Context) error {
	projectName := c.Param("projectName")

	m := model.FileSystemComposeLibraryItemUpdate{}
	r := &composeProjectUpdateRequest{ProjectName: projectName}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
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

	return ok(c, newComposeLibraryItem(m))
}