package api

import (
	"context"

	"github.com/labstack/echo/v4"

	"calendar/assets/templates"
)

func HandleIndex(c echo.Context) error {
	templates.Index().Render(context.Background(), c.Response().Writer)
	return nil
}
