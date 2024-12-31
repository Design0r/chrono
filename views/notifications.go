package views

import (
	"context"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/db/repo"
	"chrono/htmx"
	"chrono/service"
)

func InitNotificationRoutes(group *echo.Group, r *repo.Queries) {
	group.GET(
		"/notifications",
		func(c echo.Context) error { return HandleNotifications(c, r) },
	)
	group.PATCH(
		"/notifications/:id",
		func(c echo.Context) error { return HandleClearNotification(c, r) },
	)
	group.PATCH(
		"/notifications",
		func(c echo.Context) error { return HandleClearAllNotifications(c, r) },
	)
}

func HandleClearNotification(c echo.Context, r *repo.Queries) error {
	param := c.Param("id")
	num, err := strconv.Atoi(param)
	if err != nil {
		htmx.ErrorMessage("Invalid notification id", c)
		return err
	}

	currUser, err := service.GetCurrentUser(r, c)

	err = service.ClearNotification(r, int64(num))
	if err != nil {
		htmx.ErrorMessage("Failed to clear notification", c)
		return err
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		htmx.ErrorMessage("Internal Error", c)
		return err
	}

	templates.NotificationIndicator(len(notifications)).
		Render(context.Background(), c.Response().Writer)

	return nil
}

func HandleClearAllNotifications(c echo.Context, r *repo.Queries) error {
	currUser, err := service.GetCurrentUser(r, c)
	if err != nil {
		htmx.ErrorMessage("Internal error", c)
		return err
	}
	err = service.ClearAllNotifications(r, currUser.ID)
	if err != nil {
		htmx.ErrorMessage("Failed to clear notifications", c)
		return err
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		htmx.ErrorMessage("Internal Error", c)
		return err
	}

	templates.NotificationIndicator(len(notifications)).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleNotifications(c echo.Context, r *repo.Queries) error {
	currUser, err := service.GetCurrentUser(r, c)
	if err != nil {
		htmx.ErrorMessage("Internal error", c)
		return err
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		htmx.ErrorMessage("Internal Error", c)
		return err
	}

	templates.UpdateNotifications(notifications).
		Render(context.Background(), c.Response().Writer)
	return nil
}
