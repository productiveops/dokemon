package handler

import (
	"dokemon/pkg/crypto"
	"dokemon/pkg/server/model"
	"dokemon/pkg/server/requestutil"
	"errors"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func (h *Handler) CreateUser(c echo.Context) error {
	m := model.User{}
	r := &userCreateRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUniqueUsername, err := h.userStore.IsUniqueUserName(r.UserName)
	if err != nil {
		panic(err)
	}

	if !isUniqueUsername {
		return unprocessableEntity(c, duplicateUserNameError())
	}

	passwordHash, err := crypto.HashPassword(m.PasswordHash)
	if err != nil {
		panic(err)
	}
	m.PasswordHash = passwordHash

	if err := h.userStore.Create(&m); err != nil {
		panic(err)
	}

	cc := requestutil.AuthCookieContent{
		UserName: m.UserName,
		Expiry: time.Now().Add(24 * time.Hour),
	}
	requestutil.SetAuthCookie(c, cc)

	return created(c, m.Id)
}

func (h *Handler) UserLogin(c echo.Context) error {
	m := model.User{}
	r := &userLoginRequest{}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	existingUser, err := h.userStore.GetByUserName(r.UserName)
	if err != nil {
		panic(err)
	}

	if existingUser == nil {
		return unprocessableEntity(c, errors.New("invalid username or password"))
	}

	success := crypto.CheckPasswordHash(r.Password, existingUser.PasswordHash)

	if !success {
		return unprocessableEntity(c, errors.New("invalid username or password"))
	}

	cc := requestutil.AuthCookieContent{
		UserName: existingUser.UserName,
		Expiry: time.Now().Add(24 * time.Hour),
	}
	requestutil.SetAuthCookie(c, cc)

	return noContent(c)
}


func (h *Handler) UserLogout(c echo.Context) error {
	requestutil.DeleteAuthCookie(c)

	return noContent(c)
}

func (h *Handler) UserCount(c echo.Context) error {
	count, err := h.userStore.Count()
	if err != nil {
		panic(err)
	}

	return ok(c, newUserCountResponse(count))
}

func (h *Handler) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	exists, err := h.userStore.Exists(uint(id))
	if err != nil {
		panic(err)
	}

	if !exists {
		return resourceNotFound(c, "User")
	}

	m := model.User{}
	r := &userUpdateRequest{Id: uint(id)}
	if err := r.bind(c, &m); err != nil {
		return unprocessableEntity(c, err)
	}

	isUniqueUsername, err := h.userStore.IsUniqueUserName(r.UserName)
	if err != nil {
		panic(err)
	}

	if !isUniqueUsername {
		return unprocessableEntity(c, duplicateUserNameError())
	}

	if err := h.userStore.Update(&m); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) ChangeUserPassword(c echo.Context) error {
	userName := c.Get("userName").(string)
	existingUser, err := h.userStore.GetByUserName(userName)
	if err != nil {
		panic(err)
	}

	if existingUser == nil {
		return resourceNotFound(c, "User")
	}

	r := &userChangePasswordRequest{}
	if err := r.bind(c); err != nil {
		return unprocessableEntity(c, err)
	}

	correctPassword := crypto.CheckPasswordHash(r.CurrentPassword, existingUser.PasswordHash)
	if !correctPassword {
		return unprocessableEntity(c, errors.New("incorrect password"))
	}

	passwordHash, err := crypto.HashPassword(r.NewPassword)
	if err != nil {
		panic(err)
	}
	existingUser.PasswordHash = passwordHash

	if err := h.userStore.Update(existingUser); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) DeleteUserById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	exists, err := h.userStore.Exists(uint(id))
	if err != nil {
		panic(err)
	}

	if !exists {
		return resourceNotFound(c, "User")
	}

	if err := h.userStore.DeleteById(uint(id)); err != nil {
		panic(err)
	}

	return noContent(c)
}

func (h *Handler) GetUserList(c echo.Context) error {
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

	rows, totalRows, err := h.userStore.GetList(uint(p), uint(s))
	if err != nil {
		panic(err)
	}

	return ok(c, newPageResponse[userHead](newUserHeadList(rows), uint(p), uint(s), uint(totalRows)))
}

func (h *Handler) GetUserById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}

	m, err := h.userStore.GetById(uint(id))
	if err != nil {
		panic(err)
	}

	if m == nil {
		return resourceNotFound(c, "User")
	}

	return ok(c, newUserResponse(m))
}

func (h *Handler) IsUniqueUsername(c echo.Context) error {
	value := c.QueryParam("value")

	unique, err := h.userStore.IsUniqueUserName(value)
	if err != nil {
		panic(err)
	}

	return ok(c, newUniqueResponse(unique))
}

func (h *Handler) IsUniqueUsernameExcludeItself(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return unprocessableEntity(c, routeIntExpectedError("id"))
	}
	value := c.QueryParam("value")

	unique, err := h.userStore.IsUniqueUserNameExcludeItself(value, uint(id))
	if err != nil {
		panic(err)
	}

	return ok(c, newUniqueResponse(unique))
}
