package views

import (
	"context"
	"database/sql"

	"github.com/labstack/echo/v4"

	"calendar/assets/templates"
	"calendar/middleware"
	"calendar/schemas"
	"calendar/service"
)

func InitIndexRoutes(group *echo.Group, db *sql.DB) {
	group.GET(
		"",
		func(c echo.Context) error { return HandleIndex(c, db) },
		middleware.SessionMiddleware(db),
	)
}

func HandleIndex(c echo.Context, db *sql.DB) error {
	currUser, err := service.GetCurrentUser(db, c)
	if err != nil {
		return err
	}
	vacDays, err := service.GetVacationCountForUserYear(db, int(currUser.ID), service.CurrentYear())
	if err != nil {
		return err
	}

	currYear := service.CurrentYear()
	stats := schemas.YearProgress{
		NumDays:           service.NumDaysInYear(currYear),
		NumWorkDays:       service.NumWorkDays(currYear),
		NumDaysPassed:     service.YearProgress(currYear),
		DaysPassedPercent: service.YearProgressPercent(currYear),
	}

	templates.Home(currUser, vacDays, stats).Render(context.Background(), c.Response().Writer)
	return nil
}
