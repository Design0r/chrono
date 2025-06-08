package handler

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"chrono/internal/adapter/htmx"
)

func Render(c echo.Context, statusCode int, t templ.Component) error {
	return htmx.Render(c, statusCode, t)
}

func RenderError(c echo.Context, statusCode int, msg string) error {
	return htmx.RenderError(c, statusCode, msg)
}
