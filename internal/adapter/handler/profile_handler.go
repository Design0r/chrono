package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/internal/domain"
	"chrono/internal/service"
)

type ProfileHandler struct {
	user  service.UserService
	notif service.NotificationService
}

func NewProfileService(u service.UserService, n service.NotificationService) ProfileHandler {
	return ProfileHandler{user: u, notif: n}
}

func RegisterProfileRoutes(group echo.Group, handler *ProfileHandler) {
	group.GET("/profile", handler.Profile)
	group.GET("/profile/edit", handler.ProfileEditForm)
	group.PATCH("/profile", handler.)
}

func (h *ProfileHandler) Profile(c echo.Context) error {
	currUser := c.Get("user").(domain.User)
	notifications, err := h.notif.GetByUserId(c.Request().Context(), currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user notifications.")
	}

	return Render(c, http.StatusOK, templates.ProfilePage(currUser, notifications))
}


func (h *ProfileHandler) ProfileEditForm(c echo.Context) error {
	currUser := c.Get("user").(domain.User)
	notifications, err := h.notif.GetByUserId(c.Request().Context(), currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user notifications.")
	}

	return Render(c, http.StatusOK, templates.ProfileEditForm(currUser, notifications))
}


func (h *ProfileHandler) ProfileEditForm(c echo.Context) error {
	currUser := c.Get("user").(domain.User)

	patchedData, err := 

	notifications, err := h.notif.GetByUserId(c.Request().Context(), currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user notifications.")
	}

	return Render(c, http.StatusOK, templates.ProfileEditForm(currUser, notifications))
}
