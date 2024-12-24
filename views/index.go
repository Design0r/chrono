package views

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"calendar/assets/templates"
	"calendar/htmx"
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
	group.GET(
		"/notifications",
		func(c echo.Context) error { return HandleNotifications(c, db) },
		middleware.SessionMiddleware(db),
	)
	group.PATCH(
		"/notifications/:id",
		func(c echo.Context) error { return HandleClearNotification(c, db) },
		middleware.SessionMiddleware(db),
	)
	group.PATCH(
		"/notifications",
		func(c echo.Context) error { return HandleClearAllNotifications(c, db) },
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

func HandleClearNotification(c echo.Context, db *sql.DB) error {
	param := c.Param("id")
	num, err := strconv.Atoi(param)
	if err != nil {
		htmx.ErrorMessage("Invalid notification id", c)
		return err
	}

	currUser, err := service.GetCurrentUser(db, c)

	err = service.ClearNotification(db, int64(num))
	if err != nil {
		htmx.ErrorMessage("Failed to clear notification", c)
		return err
	}

	notifications, err := service.GetUserNotifications(db, currUser.ID)
	if err != nil {
		htmx.ErrorMessage("Internal Error", c)
		return err
	}

	templates.NotificationIndicator(len(notifications)).
		Render(context.Background(), c.Response().Writer)

	return nil
}

func HandleClearAllNotifications(c echo.Context, db *sql.DB) error {
	currUser, err := service.GetCurrentUser(db, c)
	if err != nil {
		htmx.ErrorMessage("Internal error", c)
		return err
	}
	err = service.ClearAllNotifications(db, currUser.ID)
	if err != nil {
		htmx.ErrorMessage("Failed to clear notifications", c)
		return err
	}

	notifications, err := service.GetUserNotifications(db, currUser.ID)
	if err != nil {
		htmx.ErrorMessage("Internal Error", c)
		return err
	}

	templates.NotificationIndicator(len(notifications)).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleNotifications(c echo.Context, db *sql.DB) error {
	currUser, err := service.GetCurrentUser(db, c)
	if err != nil {
		htmx.ErrorMessage("Internal error", c)
		return err
	}

	notifications, err := service.GetUserNotifications(db, currUser.ID)
	if err != nil {
		htmx.ErrorMessage("Internal Error", c)
		return err
	}

	templates.UpdateNotifications(notifications).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleAdminError(c echo.Context) error {
	templates.AdminError().Render(context.Background(), c.Response().Writer)
	return nil
}
