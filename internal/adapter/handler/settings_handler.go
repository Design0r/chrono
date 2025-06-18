package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/internal/adapter/htmx"
	"chrono/internal/domain"
	"chrono/internal/service"
)

type SettingsHandler struct {
	settings service.SettingsService
}

func NewSettingsHandler(s service.SettingsService) SettingsHandler {
	return SettingsHandler{settings: s}
}

func (s *SettingsHandler) RegisterRoutes(group *echo.Group) {
	g := group.Group("/settings")
	g.GET("", s.Settings)
	g.PATCH("", s.PatchSettings)
}

func (h *SettingsHandler) Settings(c echo.Context) error {
	currUser := c.Get("user").(domain.User)

	s, err := h.settings.GetFirst(c.Request().Context())
	if err != nil {
		RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.Settings(s, &currUser, []domain.Notification{}))
}

func (h *SettingsHandler) PatchSettings(c echo.Context) error {
	settings := domain.Settings{ID: 1, SignupEnabled: c.FormValue("signup_enabled") == "on"}
	_, err := h.settings.Update(c.Request().Context(), settings)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, "Failed to update settings.")
	}

	return Render(c, http.StatusOK, htmx.SuccessMessage("Successfully updated settings."))
}
