package htmx

import (
	"context"

	"github.com/labstack/echo/v4"

	"calendar/assets/templates"
)

func ErrorMessage(msg string, c echo.Context) {
	RenderMessage(msg, "error", c)
}

func InfoMessage(msg string, c echo.Context) {
	RenderMessage(msg, "info", c)
}

func SuccessMessage(msg string, c echo.Context) {
	RenderMessage(msg, "success", c)
}

func RenderMessage(msg string, mtype string, c echo.Context) {
	templates.Message(msg, mtype).Render(context.Background(), c.Response().Writer)
}

func HxRedirect(path string, c echo.Context) {
	c.Response().Header().Set("HX-Redirect", path)
}
