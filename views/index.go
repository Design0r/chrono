package views

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/htmx"
	"chrono/middleware"
	"chrono/schemas"
	"chrono/service"
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
		func(c echo.Context) error { return HandleRequests(c, db) },
		middleware.SessionMiddleware(db),
	)
	group.PATCH(
		"/requests/:id",
		func(c echo.Context) error { return HandlePatchRequests(c, db) },
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
	group.GET("/error", HandleError, middleware.SessionMiddleware(db))

	group.PUT(
		"/:id/admin",
		func(c echo.Context) error { return HandleToggleAdmin(c, db) },
		middleware.SessionMiddleware(db),
	)
}

func HandleIndex(c echo.Context, db *sql.DB) error {
	currUser, err := service.GetCurrentUser(db, c)
	if err != nil {
		htmx.ErrorPage(http.StatusNotFound, err.Error(), c)
		return err
	}
	vacDays, err := service.GetVacationCountForUserYear(db, int(currUser.ID), service.CurrentYear())
	if err != nil {
		htmx.ErrorPage(http.StatusInternalServerError, err.Error(), c)
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
		htmx.ErrorPage(http.StatusInternalServerError, err.Error(), c)
		return err
	}

	pendingEvents, err := service.GetPendingEventsForYear(db, currUser.ID, currYear)
	if err != nil {
		htmx.ErrorPage(http.StatusBadRequest, err.Error(), c)
		return err
	}

	templates.Home(currUser, vacDays, pendingEvents, stats, notifications).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleTeam(c echo.Context, db *sql.DB) error {
	currUser, err := service.GetCurrentUser(db, c)
	if err != nil {
		htmx.ErrorPage(http.StatusInternalServerError, err.Error(), c)
		return err
	}
	users, err := service.GetAllVacUsers(db)
	if err != nil {
		htmx.ErrorPage(http.StatusInternalServerError, err.Error(), c)
		return err
	}

	notifications, err := service.GetUserNotifications(db, currUser.ID)
	if err != nil {
		htmx.ErrorPage(http.StatusInternalServerError, err.Error(), c)
		return err
	}

	templates.Team(users, currUser, notifications).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleRequests(c echo.Context, db *sql.DB) error {
	currUser, err := service.GetCurrentUser(db, c)
	if err != nil {
		htmx.ErrorPage(http.StatusInternalServerError, err.Error(), c)
		return err
	}

	if !currUser.IsSuperuser {
		htmx.ErrorPage(http.StatusForbidden, "This page is only accessible by admins", c)
		return nil
	}

	requests, _ := service.GetPendingRequests(db)

	notifications, err := service.GetUserNotifications(db, currUser.ID)
	if err != nil {
		htmx.ErrorPage(http.StatusInternalServerError, err.Error(), c)
		return err
	}

	templates.Requests(&currUser, requests, notifications).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func HandlePatchRequests(c echo.Context, db *sql.DB) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		htmx.ErrorMessage("Invalid request id", c)
		return err
	}

	stateParam := c.FormValue("state")

	currUser, err := service.GetCurrentUser(db, c)
	if err != nil {
		htmx.ErrorMessage("Internal Error", c)
		return err
	}

	if !currUser.IsSuperuser {
		htmx.ErrorMessage("Not authorized", c)
		return err
	}

	err = service.UpdateRequestState(db, stateParam, currUser, int64(id))
	if err != nil {
		htmx.ErrorMessage("Failed updating request", c)
		return err
	}

	htmx.SuccessMessage(fmt.Sprintf("%v %v", strings.Title(stateParam), "Request"), c)
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

func HandleError(c echo.Context) error {
	templates.Error(http.StatusInternalServerError, "").
		Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleToggleAdmin(c echo.Context, db *sql.DB) error {
	currUser, err := service.GetCurrentUser(db, c)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	if !currUser.IsSuperuser {
		htmx.ErrorMessage("Admin rights are required to change your teams admin status", c)
		return err
	}

	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	user, err := service.ToggleAdmin(db, currUser, int64(id))
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	templates.AdminCheckbox(currUser, user).
		Render(context.Background(), c.Response().Writer)
	return nil
}
