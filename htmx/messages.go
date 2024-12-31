package htmx

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
)

func ErrorMessage(msg string, c echo.Context) templ.Component {
	return RenderMessage(msg, "error", c)
}

func InfoMessage(msg string, c echo.Context) templ.Component {
	return RenderMessage(msg, "info", c)
}

func SuccessMessage(msg string, c echo.Context) templ.Component {
	return RenderMessage(msg, "success", c)
}

func RenderMessage(msg string, mtype string, c echo.Context) templ.Component {
	return templates.Message(msg, mtype)
}

func HxRedirect(path string, c echo.Context) {
	c.Response().Header().Set("HX-Redirect", path)
}
