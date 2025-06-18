package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/internal/domain"
	"chrono/internal/service"
)

type NotificationHandler struct {
	log   *slog.Logger
	notif service.NotificationService
}

func NewNotificationHandler(n service.NotificationService, log *slog.Logger) NotificationHandler {
	return NotificationHandler{notif: n, log: log}
}

func (h *NotificationHandler) RegisterRoutes(group *echo.Group) {
	group.GET("/notifications", h.Notifications)
	group.PATCH("/notifications/:id", h.ClearNotification)
	group.PATCH("/notifications", h.ClearAllNotifications)
}

func (h *NotificationHandler) Notifications(c echo.Context) error {
	currUser := c.Get("user").(domain.User)

	notifications, err := h.notif.GetByUserId(c.Request().Context(), currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user notifications.")
	}

	return Render(c, http.StatusOK, templates.UpdateNotifications(notifications))
}

func (h *NotificationHandler) ClearNotification(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	param := c.Param("id")
	notifId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid notification id")
	}

	err = h.notif.Clear(ctx, notifId)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to clear notification")
	}

	notifications, err := h.notif.GetByUserId(ctx, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.NotificationIndicator(len(notifications)))
}

func (h *NotificationHandler) ClearAllNotifications(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	err := h.notif.ClearAll(ctx, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to clear notifications.")
	}

	notifications, err := h.notif.GetByUserId(ctx, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user notifications.")
	}

	return Render(c, http.StatusOK, templates.NotificationIndicator(len(notifications)))
}
