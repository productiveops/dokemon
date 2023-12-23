package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

// General

type pageResponse[T any] struct {
	Items []T 			`json:"items"`
	PageNo uint 		`json:"pageNo"`
	PageSize uint 		`json:"pageSize"`
	TotalRows uint 		`json:"totalRows"`
}

type entityCreatedResponse struct {
	Id any 		`json:"id"`
}

type uniqueResponse struct {
	Unique bool `json:"unique"`
}

func newPageResponse[T any](rows []T, pageNo, pageSize, totalRows uint) *pageResponse[T] {
	return &pageResponse[T]{
		Items: rows,
		PageNo: pageNo,
		PageSize: pageSize,
		TotalRows: totalRows,
	}
}

func newUniqueResponse(unique bool) *uniqueResponse {
	return &uniqueResponse{
		Unique: unique,
	}
}

func ok(c echo.Context, r interface{}) error {
	return c.JSON(http.StatusOK, r)
}

func created(c echo.Context, id any) error {
	return c.JSON(http.StatusCreated, entityCreatedResponse{Id: id})
}

func noContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func unprocessableEntity(c echo.Context, err error) error {
	return c.JSON(http.StatusUnprocessableEntity, newError(err))
}

func resourceNotFound(c echo.Context, resourceName string) error {
	return c.JSON(http.StatusNotFound, resourceNotFoundError(resourceName))
}

func notFound(c echo.Context, message string) error {
	return c.JSON(http.StatusNotFound, errors.New(message))
}
