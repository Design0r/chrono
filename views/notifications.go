package views

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/db/repo"
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
	currUser := c.Get("user").(repo.User)

	param := c.Param("id")
	num, err := strconv.Atoi(param)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid notification id")
	}

	err = service.ClearNotification(r, int64(num))
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to clear notification")
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.NotificationIndicator(len(notifications)))
}

func HandleClearAllNotifications(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)

	err := service.ClearAllNotifications(r, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to clear notification")
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.NotificationIndicator(len(notifications)))
}

func HandleNotifications(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.UpdateNotifications(notifications))
}
