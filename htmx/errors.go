package htmx

import (
	"context"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
)

func ErrorPage(statusCode int, message string, c echo.Context) {
	templates.Error(statusCode, message).Render(context.Background(), c.Response().Writer)
}
