package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/internal/domain"
	"chrono/internal/service"
)

type DebugHandler struct {
	user    service.UserService
	notif   service.NotificationService
	token   service.TokenService
	session service.SessionService
	log     *slog.Logger
}

func NewDebugHandler(
	u service.UserService,
	n service.NotificationService,
	t service.TokenService,
	s service.SessionService,
	log *slog.Logger,
) DebugHandler {
	return DebugHandler{user: u, notif: n, token: t, session: s, log: log}
}

func (h *DebugHandler) RegisterRoutes(group *echo.Group) {
	group.GET("/debug", h.Debug)
	group.DELETE("/debug/tokens", h.DeleteTokens)
	group.DELETE("/debug/sessions", h.DeleteSessions)
	group.PATCH("/debug/color", h.UserColor)
}

func (h *DebugHandler) Debug(c echo.Context) error {
	currUser := c.Get("user").(domain.User)

	notifications, err := h.notif.GetByUserId(c.Request().Context(), currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user notifications.")
	}

	return Render(
		c,
		http.StatusOK,
		templates.Debug(&currUser, notifications),
	)
}

func (h *DebugHandler) DeleteTokens(c echo.Context) error {
	err := h.token.DeleteAll(c.Request().Context())
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.Message("Reset token table", "success"))
}

func (h *DebugHandler) DeleteSessions(c echo.Context) error {
	err := h.session.DeleteAll(c.Request().Context())
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.Message("Reset session table", "success"))
}

func (h *DebugHandler) UserColor(c echo.Context) error {
	ctx := c.Request().Context()
	users, err := h.user.GetAll(ctx)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	for _, user := range users {
		user.Color = domain.Color.HSLToHex(domain.Color.HSLFloat(1))
		h.user.Update(ctx, &user)
	}

	bot, err := h.user.GetById(ctx, 1)
	bot.Color = domain.Color.HSLToHex(domain.Color.HSLFloat(1))
	h.user.Update(ctx, bot)

	return Render(c, http.StatusOK, templates.Message("Changed user default colors", "success"))
}
