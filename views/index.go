package views

import (
	"context"
	"database/sql"

	"github.com/labstack/echo/v4"

	"calendar/assets/templates"
	"calendar/middleware"
)

func InitIndexRoutes(group *echo.Group, db *sql.DB) {
	group.GET(
		"",
		func(c echo.Context) error { return HandleIndex(c, db) },
		middleware.SessionMiddleware(db),
	)
}

func HandleIndex(c echo.Context, db *sql.DB) error {
	templates.Home().Render(context.Background(), c.Response().Writer)
	return nil
}
