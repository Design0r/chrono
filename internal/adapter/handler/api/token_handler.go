package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"chrono/internal/domain"
	"chrono/internal/service"
)

type APITokenHandler struct {
	user  *service.UserService
	token *service.TokenService
	notif *service.NotificationService
	log   *slog.Logger
}

func NewAPITokenHandler(
	t *service.TokenService,
	u *service.UserService,
	n *service.NotificationService,
	log *slog.Logger,
) APITokenHandler {
	return APITokenHandler{token: t, user: u, notif: n, log: log}
}

func (h *APITokenHandler) RegisterRoutes(group *echo.Group) {
	group.POST("/tokens", h.CreateTokens)
}

func (h *APITokenHandler) CreateTokens(c echo.Context) error {
	ctx := c.Request().Context()
	currUser := c.Get("user").(domain.User)
	params, err := c.FormParams()
	if err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	userIdParam := params["filter"][0]
	userId, err := strconv.ParseInt(userIdParam, 10, 64)
	if err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	tokenValue := params["token"][0]
	tokenNum, err := strconv.ParseFloat(tokenValue, 64)
	if err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	user, err := h.user.GetById(ctx, userId)
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, "Failed to get user.")
	}

	_, err = h.token.CreateVacationToken(ctx, tokenNum, time.Now().Year(), user.ID)
	if err != nil {
		return NewErrorResponse(
			c,
			http.StatusInternalServerError,
			"Failed to create vacation token.",
		)
	}

	msg := fmt.Sprintf("You received %v vacation token from %v", tokenNum, currUser.Username)
	err = h.notif.CreateAndNotify(ctx, msg, []domain.User{*user})
	if err != nil {
		return NewErrorResponse(c, http.StatusInternalServerError, "Failed to create notification.")
	}

	return NewJsonResponse(c, nil)
}
