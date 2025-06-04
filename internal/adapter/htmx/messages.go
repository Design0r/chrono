package htmx

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
)

func ErrorMessage(msg string) templ.Component {
	return RenderMessage(msg, "error")
}

func InfoMessage(msg string) templ.Component {
	return RenderMessage(msg, "info")
}

func SuccessMessage(msg string) templ.Component {
	return RenderMessage(msg, "success")
}

func RenderMessage(msg string, mtype string) templ.Component {
	return templates.Message(msg, mtype)
}

func HxRedirect(path string, c echo.Context) {
	c.Response().Header().Set("HX-Redirect", path)
}
