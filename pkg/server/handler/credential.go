package handler

import (
	"errors"
	"strconv"

	"github.com/productiveops/dokemon/pkg/server/model"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateCredential(c echo.Context) error {
	m := model.Credential{}
	r := &credentialCreateRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUnique, err := h.credentialStore.IsUniqueName(r.Name)
	if err != nil {
		panic(err)
	}

	if !isUnique {
		return unprocessableEntity(c, duplicateNameError())
	}

	if err := h.credentialStore.Create(&m); err != nil {
		panic(err)
	}

	return created(c, m.Id)
}

func (h *Handler) UpdateCredentialDetails(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}
	
	m, err := h.credentialStore.GetById(uint(id))
	if err != nil {
		panic(err)
	}
	
	if m == nil {
		return resourceNotFound(c, "Credential")
	}	

	r := &credentialUpdateDetailsRequest{Id: uint(id)}
	if err := r.bind(c, m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUnique, err := h.credentialStore.IsUniqueNameExcludeItself(r.Name, r.Id)
	if err != nil {
		panic(err)
	}

	if !isUnique {
		return unprocessableEntity(c, duplicateNameError())
	}

	if err := h.credentialStore.Update(m); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) UpdateCredentialSecret(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}
	
	m, err := h.credentialStore.GetById(uint(id))
	if err != nil {
		panic(err)
	}
	
	if m == nil {
		return resourceNotFound(c, "Credential")
	}	

	r := &credentialUpdateSecretRequest{Id: uint(id)}
	if err := r.bind(c, m); err != nil {
		return unprocessableEntity(c, err)
	}

	if err := h.credentialStore.Update(m); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) DeleteCredentialById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	exists, err := h.credentialStore.Exists(uint(id))
	if err != nil {
		panic(err)
	}

	if !exists {
		return resourceNotFound(c, "Credential")
	}

	inUse, err := h.credentialStore.IsInUse(uint(id))
	if err != nil {
		panic(err)
	}

	if inUse {
		return unprocessableEntity(c, errors.New("Credentials are in use and cannot be deleted"))
	}

	if err := h.credentialStore.DeleteById(uint(id)); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) GetCredentialList(c echo.Context) error {
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

	rows, totalRows, err := h.credentialStore.GetList(uint(p), uint(s))
	if err != nil {
		panic(err)
	}

	return ok(c, newPageResponse[credentialHead](newCredentialHeadList(rows), uint(p), uint(s), uint(totalRows)))
}

func (h *Handler) GetCredentialById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	m, err := h.credentialStore.GetById(uint(id))
	if err != nil {
		panic(err)
	}

	if m == nil {
		return resourceNotFound(c, "Credential")
	}

	return ok(c, newCredentialResponse(m))
}

func (h *Handler) IsUniqueCredentialName(c echo.Context) error {
	value := c.QueryParam("value")

	unique, err := h.credentialStore.IsUniqueName(value)
	if err != nil {
		panic(err)
	}

	return ok(c, newUniqueResponse(unique))
}

func (h *Handler) IsUniqueCredentialNameExcludeItself(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	value := c.QueryParam("value")

	unique, err := h.credentialStore.IsUniqueNameExcludeItself(value, uint(id))
	if err != nil {
		panic(err)
	}

	return ok(c, newUniqueResponse(unique))
}
