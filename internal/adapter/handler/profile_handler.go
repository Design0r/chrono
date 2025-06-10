package handler

import (
	"net/http"
	"strconv"

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
	group.PATCH("/profile", handler.ProfileEdit)
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

func (h *ProfileHandler) ProfileEdit(c echo.Context) error {
	currUser := c.Get("user").(domain.User)

	patchedData := domain.PatchUser{}
	if err := c.Bind(&patchedData); err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid data")
	}

	updatedUser, err := h.user.Update(
		c.Request().Context(),
		&domain.User{
			ID:           currUser.ID,
			Username:     patchedData.Name,
			Email:        patchedData.Email,
			Color:        patchedData.Color,
			IsSuperuser:  currUser.IsSuperuser,
			VacationDays: currUser.VacationDays,
		},
	)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to update user information.")
	}

	notifications, err := h.notif.GetByUserId(c.Request().Context(), currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user notifications.")
	}

	return Render(c, http.StatusOK, templates.UpdateProfileWithMessage(*updatedUser, notifications))
}

func (h *ProfileHandler) ToggleAdmin(c echo.Context) error {
	currUser := c.Get("user").(domain.User)

	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid user id")
	}

	updatedUser, err := h.user.ToggleAdmin(c.Request().Context(), id, &currUser)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to change admin status")
	}

	return Render(c, http.StatusOK, templates.AdminCheckbox(currUser, updatedUser.ID, updatedUser.IsSuperuser, true))
}
