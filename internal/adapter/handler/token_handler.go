package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/internal/domain"
	"chrono/internal/service"
)

type TokenHandler struct {
	user     service.UserService
	vacation service.VacationTokenService
	notif    service.NotificationService
	log      *slog.Logger
}

func NewTokenHandler(
	v service.VacationTokenService,
	u service.UserService,
	n service.NotificationService,
	log *slog.Logger,
) TokenHandler {
	return TokenHandler{vacation: v, user: u, notif: n, log: log}
}

func (h *TokenHandler) RegisterRoutes(group *echo.Group) {
	group.GET("/tokens", h.Token)
	group.POST("/tokens", h.CreateTokens)
}

func (h *TokenHandler) Token(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)

	users, err := h.user.GetAll(ctx)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	notifications, err := h.notif.GetByUserId(ctx, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.Tokens(&currUser, notifications, users))
}

func (h *TokenHandler) CreateTokens(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)
	params, err := c.FormParams()
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	userName := params["filter"][0]
	tokenValue := params["token"][0]

	tokenNum, err := strconv.ParseFloat(tokenValue, 64)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	user, err := h.user.GetByName(ctx, userName)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to get user.")
	}

	_, err = h.vacation.Create(ctx, tokenNum, time.Now().Year(), user.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to create vacation token.")
	}

	msg := fmt.Sprintf("You received %v vacation token from %v", tokenNum, currUser.Username)
	err = h.notif.CreateAndNotify(ctx, msg, []domain.User{*user})
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, "Failed to create notification.")
	}

	return Render(c, http.StatusOK, templates.Message("Created Token", "success"))
}
