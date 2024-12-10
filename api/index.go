package api

import (
	"database/sql"

	"github.com/labstack/echo/v4"
)

func InitIndexRoutes(group *echo.Group, db *sql.DB) {
	group.GET("", func(c echo.Context) error { return HandleIndex(c, db) })
}

func HandleIndex(c echo.Context, db *sql.DB) error {
	// templates.Index(events).Render(context.Background(), c.Response().Writer)
	return nil
}
