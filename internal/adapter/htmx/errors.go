package htmx

import (
	"context"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
)

func ErrorPage(statusCode int, message string, c echo.Context) {
	templates.Error(statusCode, message).Render(context.Background(), c.Response().Writer)
}

func IsHTMXRequest(c echo.Context) bool {
	if _, exists := c.Request().Header["HX-Request"]; exists {
		return true
	}
	if _, exists := c.Request().Header["Hx-Request"]; exists {
		return true
	}

	return false
}
