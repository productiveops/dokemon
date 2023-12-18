package handler

import (
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

	rows, totalRows, err := h.localComposeLibraryStore.GetList(uint(p), uint(s))
	if err != nil {
		panic(err)
	}

	return ok(c, newPageResponse[composeLibraryItemHead](newComposeLibraryItemHeadList(rows), uint(p), uint(s), uint(totalRows)))
}

func (h *Handler) CreateComposeProject(c echo.Context) error {
	m := model.LocalComposeLibraryItem{}
	r := &composeProjectCreateRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	if err := h.localComposeLibraryStore.Create(&m); err != nil {
		panic(err)
	}

	return created(c, m.ProjectName)
}

func (h *Handler) UpdateComposeProject(c echo.Context) error {
	projectName := c.Param("projectName")

	m := model.LocalComposeLibraryItemUpdate{}
	r := &composeProjectUpdateRequest{ProjectName: projectName}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	if err := h.localComposeLibraryStore.Update(&m); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) DeleteComposeProject(c echo.Context) error {
	projectName := c.Param("projectName")

	if err := h.localComposeLibraryStore.DeleteByName(projectName); err != nil {
		return unprocessableEntity(c, err)
	}

	return noContent(c)
}

func (h *Handler) GetComposeProject(c echo.Context) error {
	projectName := c.Param("projectName")

	m, err := h.localComposeLibraryStore.GetByName(projectName)
	if err != nil {
		return notFound(c, "ComposeProject")
	}

	return ok(c, newComposeLibraryItem(m))
}