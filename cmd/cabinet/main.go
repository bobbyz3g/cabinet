package main

import (
	"net/http"

	"github.com/bobbyz3g/cabinet"
	"github.com/labstack/echo/v4"
)

func main() {
	sessions := &cabinet.Sessions{}

	e := echo.New()
	e.POST("/file", func(c echo.Context) error {
		h := &cabinet.PushHandler{Sessions: sessions}
		h.Prepare(c)
		if err := h.Flush(); err != nil {
			return err
		}
		return c.String(http.StatusCreated, "ok")
	})
	e.Logger.Fatal(e.Start(":8080"))
}
