package api

import (
	"context"
	"database/sql"
	"time"

	"github.com/labstack/echo/v4"

	"calendar/assets/templates"
	"calendar/service"
)

func InitIndexRoutes(group *echo.Group, db *sql.DB) {
	group.GET("", func(c echo.Context) error { return HandleIndex(c, db) })
}

func HandleIndex(c echo.Context, db *sql.DB) error {
	events, err := service.GetEventsForMonth(
		db,
		time.Now(),
	)
	if err != nil {
		return err
	}
	templates.Index(events).Render(context.Background(), c.Response().Writer)
	return nil
}
