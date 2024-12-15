package views

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"calendar/assets/templates"
	"calendar/middleware"
	"calendar/schemas"
	"calendar/service"
)

func InitEventRoutes(group *echo.Group, db *sql.DB) {
	group.POST(
		"/:year/:month/:day",
		func(c echo.Context) error { return CreateEventHandler(c, db) },
		middleware.SessionMiddleware(db),
	)
	group.DELETE(
		"/events/:id",
		func(c echo.Context) error { return DeleteEventHandler(c, db) },
		middleware.SessionMiddleware(db),
	)
}

func CreateEventHandler(c echo.Context, db *sql.DB) error {
	var date schemas.YMDDate
	if err := c.Bind(&date); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid parameter")
	}
	eventName := c.FormValue("name")

	currUser, err := service.GetCurrentUser(db, c)
	if err := c.Bind(&date); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Current user not found")
	}

	event, err := service.CreateEvent(db, date, currUser.ID, eventName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}

	templates.Event(schemas.Event{Username: currUser.Username, Event: event}).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func DeleteEventHandler(c echo.Context, db *sql.DB) error {
	event := c.Param("id")
	eventId, err := strconv.Atoi(event)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid event id")
	}

	err = service.DeleteEvent(db, eventId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid event id")
	}
	return nil
}
