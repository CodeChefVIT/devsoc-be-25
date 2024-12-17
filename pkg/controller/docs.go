package controller

import (
	"net/http"

	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/labstack/echo/v4"
)

func Docs(c echo.Context) error {
	content, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: "./docs/sample.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "DevSoc 2025",
		},
		DarkMode: true,
	})

	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, content)
}
