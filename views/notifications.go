package views

import (
	"net/http"
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
		return Render(c, http.StatusBadRequest, htmx.ErrorMessage("Invalid notification id", c))
	}

	currUser, err := service.GetCurrentUser(r, c)

	err = service.ClearNotification(r, int64(num))
	if err != nil {
		return Render(
			c,
			http.StatusInternalServerError,
			htmx.ErrorMessage("Failed to clear notification", c),
		)
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return Render(
			c,
			http.StatusInternalServerError,
			htmx.ErrorMessage("Internal server error", c),
		)
	}

	return Render(c, http.StatusOK, templates.NotificationIndicator(len(notifications)))
}

func HandleClearAllNotifications(c echo.Context, r *repo.Queries) error {
	currUser, err := service.GetCurrentUser(r, c)
	if err != nil {
		return Render(
			c,
			http.StatusInternalServerError,
			htmx.ErrorMessage("Internal server error", c),
		)
	}
	err = service.ClearAllNotifications(r, currUser.ID)
	if err != nil {
		return Render(
			c,
			http.StatusInternalServerError,
			htmx.ErrorMessage("Failed to clear notifications", c),
		)
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return Render(
			c,
			http.StatusInternalServerError,
			htmx.ErrorMessage("Internal server error", c),
		)
	}

	return Render(c, http.StatusOK, templates.NotificationIndicator(len(notifications)))
}

func HandleNotifications(c echo.Context, r *repo.Queries) error {
	currUser, err := service.GetCurrentUser(r, c)
	if err != nil {
		return Render(
			c,
			http.StatusInternalServerError,
			htmx.ErrorMessage("Internal server error", c),
		)
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return Render(
			c,
			http.StatusInternalServerError,
			htmx.ErrorMessage("Internal server error", c),
		)
	}

	return Render(c, http.StatusOK, templates.UpdateNotifications(notifications))
}
