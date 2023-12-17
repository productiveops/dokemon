package web

import (
	"io"
	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo) {
    // Ref: https://anhduongviet.medium.com/combine-go-and-react-in-single-docker-container-28e4df0c2d48
    //      https://github.com/vietanhduong/go-n-reactjs
    frontend := rice.MustFindBox("dist")
    fe := http.FileServer(frontend.HTTPBox())

    e.GET("/assets/*", echo.WrapHandler(fe))

    e.GET("/*", func(c echo.Context) error {
      index, err := frontend.Open("index.html")
      if err != nil {
         return err
      }
      content, err := io.ReadAll(index)
      if err != nil {
         return err
      }
      return c.HTMLBlob(http.StatusOK, content)
   })
}