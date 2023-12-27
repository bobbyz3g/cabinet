package main

import (
	"io"
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

	e.GET("/file/:code", func(c echo.Context) error {
		code := cabinet.Code(c.Param("code"))
		t, ok := sessions.Pop(code)
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound, "translator not found")
		}
		defer close(t.Done)

		c.Response().Header().Set("Content-Disposition", "attachment; filename="+t.Name)
		c.Response().Header().Set("Content-Type", "application/octet-stream")
		_, err := io.Copy(c.Response(), t.Reader)
		if err != nil {
			return err
		}
		return nil
	})
	e.Logger.Fatal(e.Start(":8080"))
}
