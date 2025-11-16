package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/internal/domain"
	"chrono/internal/service"
)

type APIAuthHandler struct {
	user service.UserService
	auth service.AuthService
	log  *slog.Logger
}

func NewAPIAuthHandler(
	u service.UserService,
	a service.AuthService,
	log *slog.Logger,
) APIAuthHandler {
	return APIAuthHandler{user: u, auth: a, log: log}
}

func (h *APIAuthHandler) RegisterRoutes(group *echo.Group) {
	group.POST("/login", h.Login)
	group.POST("/signup", h.Signup)
	group.POST("/logout", h.Logout)
}

func (h *APIAuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()

	var loginData domain.APILogin
	if err := c.Bind(&loginData); err != nil {
		return NewErrorResponse(c, http.StatusBadRequest, "Invalid inputs")
	}

	fmt.Println(loginData)

	cookie, err := h.auth.Login(ctx, loginData.Email, loginData.Password)
	if err != nil {
		return NewErrorResponse(c, http.StatusNotFound, "Incorrect email or password")
	}

	user, err := h.user.GetByEmail(ctx, loginData.Email)
	if err != nil {
		return NewErrorResponse(c, http.StatusNotFound, "Unable to fetch user")
	}

	c.SetCookie(cookie)

	return NewJsonResponse(c, user)
}

func (h *APIAuthHandler) Signup(c echo.Context) error {
	ctx := c.Request().Context()
	settings := c.Get("settings").(domain.Settings)
	if !settings.SignupEnabled {
		return NewErrorResponse(c, http.StatusForbidden, "Signups are currently disabled.")
	}

	var loginData domain.ApiCreateUser
	if err := c.Bind(&loginData); err != nil {
		return NewErrorResponse(c, http.StatusUnprocessableEntity, "Invalid inputs")
	}

	cookie, err := h.auth.Signup(
		ctx,
		domain.CreateUser{
			Username: loginData.Username,
			Password: loginData.Password,
			Email:    loginData.Email,
		},
	)
	if err != nil {
		return NewErrorResponse(c, http.StatusNotFound, err.Error())
	}

	user, err := h.user.GetByEmail(ctx, loginData.Email)
	if err != nil {
		return NewErrorResponse(c, http.StatusNotFound, err.Error())
	}

	c.SetCookie(cookie)
	return NewJsonResponse(c, user)
}

func (h *APIAuthHandler) Logout(c echo.Context) error {
	session, err := c.Cookie("session")
	if err != nil {
		return NewErrorResponse(
			c,
			http.StatusBadRequest,
			"No active user session found. Already logged out",
		)
	}
	cookie, err := h.auth.Logout(c.Request().Context(), session.Value)
	if err != nil {
		return NewErrorResponse(
			c,
			http.StatusInternalServerError,
			"Logout failed",
		)
	}

	c.SetCookie(cookie)
	return NewJsonResponse(c, nil)
}
