package htmx

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func RenderError(c echo.Context, statusCode int, msg string) error {
	if IsHTMXRequest(c) {
		return Render(c, http.StatusOK, ErrorMessage(msg))
	}
	return Render(c, statusCode, templates.Error(statusCode, msg))
}
