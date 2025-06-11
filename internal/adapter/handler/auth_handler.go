package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/internal/adapter/htmx"
	"chrono/internal/domain"
	"chrono/internal/service"
)

type AuthHandler struct {
	user service.UserService
	auth service.AuthService
	log  *slog.Logger
}

func NewAuthHandler(u service.UserService, a service.AuthService, log *slog.Logger) AuthHandler {
	return AuthHandler{user: u, auth: a, log: log}
}

func (h *AuthHandler) RegisterRoutes(group *echo.Group) {
	group.GET("/login", h.LoginForm)
	group.GET("/signup", h.SignupForm)

	group.POST("/login", h.Login)
	group.POST("/signup", h.Signup)
	group.POST("/logout", h.Logout)
}

func (h *AuthHandler) LoginForm(c echo.Context) error {
	return Render(c, http.StatusOK, templates.Login())
}

func (h *AuthHandler) Login(c echo.Context) error {
	var loginData domain.Login
	if err := c.Bind(&loginData); err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid inputs")
	}

	cookie, err := h.auth.Login(c.Request().Context(), loginData.Email, loginData.Password)
	if err != nil {
		return RenderError(c, http.StatusNotFound, "Incorrect email or password")
	}

	c.SetCookie(cookie)
	htmx.HxRedirect("/", c)
	return nil
}

func (h *AuthHandler) SignupForm(c echo.Context) error {
	return Render(c, http.StatusOK, templates.Signup())
}

func (h *AuthHandler) Signup(c echo.Context) error {
	var loginData domain.Login
	if err := c.Bind(&loginData); err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid inputs")
	}

	cookie, err := h.auth.Login(c.Request().Context(), loginData.Email, loginData.Password)
	if err != nil {
		return RenderError(c, http.StatusNotFound, "Incorrect email or password")
	}

	c.SetCookie(cookie)
	htmx.HxRedirect("/", c)
	return nil
}

func (h *AuthHandler) Logout(c echo.Context) error {
	session, err := c.Cookie("session")
	if err != nil {
		return RenderError(
			c,
			http.StatusBadRequest,
			"No active user session found. Already logged out",
		)
	}
	cookie, err := h.auth.Logout(c.Request().Context(), session.Value)
	if err != nil {
		return RenderError(
			c,
			http.StatusInternalServerError,
			"Logout failed",
		)
	}

	c.SetCookie(cookie)
	return c.Redirect(http.StatusFound, "/auth/login")
}
