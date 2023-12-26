package main

import (
	"net/http"
	"os"

	"github.com/bobbyz3g/cabinet"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.POST("/file", func(c echo.Context) error {
		code := cabinet.GenerateCode()
		h := &cabinet.SaveHandler{Ctx: c}
		h.FormFile("file")
		h.OpenFileHeader()
		h.OpenOSFile(code, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		h.Save()
		if err := h.Err(); err != nil {
			return err
		}
		return c.String(http.StatusCreated, code)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
