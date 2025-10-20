package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/config"
	"chrono/internal/domain"
	"chrono/internal/service"
)

type DebugHandler struct {
	user    service.UserService
	auth    service.AuthService
	notif   service.NotificationService
	token   service.TokenService
	session service.SessionService
	log     *slog.Logger
	event   service.EventService
}

func NewDebugHandler(
	u service.UserService,
	a service.AuthService,
	n service.NotificationService,
	t service.TokenService,
	s service.SessionService,
	e service.EventService,
	log *slog.Logger,
) DebugHandler {
	return DebugHandler{user: u, auth: a, notif: n, token: t, session: s, log: log, event: e}
}

func (h *DebugHandler) RegisterRoutes(group *echo.Group) {
	group.GET("/debug", h.Debug)
	group.DELETE("/debug/tokens", h.DeleteTokens)
	group.DELETE("/debug/sessions", h.DeleteSessions)
	group.PATCH("/debug/color", h.UserColor)
	group.PATCH("/debug/password", h.ChangePassword)
	group.DELETE("/debug/botEvents", h.DeleteBotEventsByName)
}

func (h *DebugHandler) Debug(c echo.Context) error {
	currUser := c.Get("user").(domain.User)
	ctx := c.Request().Context()

	notifications, err := h.notif.GetByUserId(c.Request().Context(), currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user notifications.")
	}

	users, err := h.user.GetAll(ctx)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get users.")
	}

	return Render(
		c,
		http.StatusOK,
		templates.Debug(&currUser, users, notifications),
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
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}
	bot.Color = domain.Color.HSLToHex(domain.Color.HSLFloat(1))
	h.user.Update(ctx, bot)

	return Render(c, http.StatusOK, templates.Message("Changed user default colors", "success"))
}

func (h *DebugHandler) ChangePassword(c echo.Context) error {
	ctx := c.Request().Context()
	userName := c.FormValue("user")
	newPw := c.FormValue("password")

	if userName == "" && newPw == "" {
		return RenderError(c, http.StatusBadRequest, "Username and password cant be empty")
	}

	user, err := h.user.GetByName(ctx, userName)
	if err != nil {
		return RenderError(c, http.StatusNotFound, "User not found")
	}
	pw, err := h.auth.HashPassword(newPw)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Unable to hash password")
	}

	user.Password = pw

	_, err = h.user.Update(ctx, user)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Unable to hash password")
	}

	return Render(c, http.StatusOK, templates.Message("Successfully changed password", "success"))
}

func (h *DebugHandler) DeleteBotEventsByName(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)
	eventName := c.FormValue("eventName")
	cfg := config.GetConfig()
	bot, err := h.user.GetByName(ctx, cfg.BotName)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get chrono bot.")
	}
	events, err := h.event.GetAllByUserId(ctx, bot.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get events for user.")
	}

	for _, event := range events {
		if event.Name == eventName {
			_, err = h.event.Delete(ctx, event.ID, &currUser)
			if err != nil {
				h.log.Error("failed deleting event", slog.String("event", event.Name), slog.String("error", err.Error()))
			}

			h.log.Debug("deleted event", slog.String("event", event.Name))
		}
	}

	return Render(c, http.StatusOK, templates.Message("Successfully deleted events", "success"))
}
