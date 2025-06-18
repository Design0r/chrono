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

type ProfileHandler struct {
	user  service.UserService
	notif service.NotificationService
	log   *slog.Logger
}

func NewProfileHandler(
	u service.UserService,
	n service.NotificationService,
	log *slog.Logger,
) ProfileHandler {
	return ProfileHandler{user: u, notif: n, log: log}
}

func (h *ProfileHandler) RegisterRoutes(authGrp *echo.Group, adminGrp *echo.Group) {
	authGrp.GET("/profile", h.Profile)
	authGrp.GET("/profile/edit", h.ProfileEditForm)
	authGrp.PATCH("/profile", h.ProfileEdit)
	adminGrp.PATCH("/profile/:id/admin", h.ToggleAdmin)
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

	return Render(
		c,
		http.StatusOK,
		templates.AdminCheckbox(currUser, updatedUser.ID, updatedUser.IsSuperuser, true),
	)
}
