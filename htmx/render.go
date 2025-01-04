package htmx

import (
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
)

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	// csrf := ctx.Get("csrf").(string)
	csrf := "1234354"
	reqCtx := context.WithValue(ctx.Request().Context(), "csrf", csrf)

	if err := t.Render(reqCtx, buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}

func RenderError(c echo.Context, statusCode int, msg string) error {
	if IsHTMXRequest(c) {
		return Render(c, statusCode, ErrorMessage(msg))
	}
	return Render(c, statusCode, templates.Error(statusCode, msg))
}
