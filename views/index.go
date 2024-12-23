package views

import (
	"context"
	"database/sql"
	"net/http"

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
	group.GET(
		"/team",
		func(c echo.Context) error { return HandleTeam(c, db) },
		middleware.SessionMiddleware(db),
	)
	group.GET(
		"/requests",
		func(c echo.Context) error { return HandleTeam(c, db) },
		middleware.SessionMiddleware(db),
	)
	group.GET("/errors/admin", HandleAdminError, middleware.SessionMiddleware(db))
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

	notifications, err := service.GetUserNotifications(db, currUser.ID)
	if err != nil {
		return err
	}

	templates.Home(currUser, vacDays, stats, notifications).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleTeam(c echo.Context, db *sql.DB) error {
	currUser, err := service.GetCurrentUser(db, c)
	if err != nil {
		return err
	}
	users, err := service.GetAllVacUsers(db)
	if err != nil {
		return err
	}

	notifications, err := service.GetUserNotifications(db, currUser.ID)
	if err != nil {
		return err
	}

	templates.Team(users, currUser, notifications).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleRequests(c echo.Context, db *sql.DB) error {
	currUser, err := service.GetCurrentUser(db, c)
	if err != nil {
		return err
	}

	if !currUser.IsSuperuser {
		return c.Redirect(http.StatusFound, "/admin-error")
	}

	requests, err := service.GetPendingRequests(db)
	if err != nil {
		return err
	}

	notifications, err := service.GetUserNotifications(db, currUser.ID)
	if err != nil {
		return err
	}

	templates.Requests(&currUser, requests, notifications).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleAdminError(c echo.Context) error {
	templates.AdminError().Render(context.Background(), c.Response().Writer)
	return nil
}
