package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"

	"calendar/assets/templates"
	"calendar/schemas"
	"calendar/service"
)

func InitEventRoutes(group *echo.Group, db *sql.DB) {
	group.POST(
		"/:year/:month/:day",
		func(c echo.Context) error { return CreateEventHandler(c, db) },
	)
}

func CreateEventHandler(c echo.Context, db *sql.DB) error {
	var date schemas.Date
	if err := c.Bind(&date); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid parameter")
	}

	event, err := service.CreateEvent(db, date)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}

	templates.Event(event).
		Render(context.Background(), c.Response().Writer)
	return nil
}
