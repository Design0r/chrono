package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"chrono/internal/domain"
	"chrono/internal/service"
)

type APINotificationHandler struct {
	log   *slog.Logger
	notif service.NotificationService
}

func NewAPINotificationHandler(n service.NotificationService, log *slog.Logger) APINotificationHandler {
	return APINotificationHandler{notif: n, log: log}
}

func (h *APINotificationHandler) RegisterRoutes(group *echo.Group) {
	group.GET("/notifications", h.Notifications)
	group.PATCH("/notifications/:id", h.ClearNotification)
	group.PATCH("/notifications", h.ClearAllNotifications)
}

func (h *APINotificationHandler) Notifications(c echo.Context) error {
	currUser := c.Get("user").(domain.User)

	notifications, err := h.notif.GetByUserId(c.Request().Context(), currUser.ID)
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, "Failed to get user notifications.")
	}

	return NewJsonResponse(c, notifications)
}

func (h *APINotificationHandler) ClearNotification(c echo.Context) error {
	ctx := c.Request().Context()

	param := c.Param("id")
	notifId, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "Invalid notification id")
	}

	err = h.notif.Clear(ctx, notifId)
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, "Failed to clear notification")
	}

	return NewJsonResponse(c, nil)
}

func (h *APINotificationHandler) ClearAllNotifications(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	err := h.notif.ClearAll(ctx, currUser.ID)
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, "Failed to clear notifications.")
	}

	return NewJsonResponse(c, nil)
}
