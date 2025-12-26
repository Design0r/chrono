package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"chrono/internal/adapter/handler/api"
	"chrono/internal/domain"
	"chrono/internal/service"
)

type MiddlewareFunc = func(echo.HandlerFunc) echo.HandlerFunc

func SessionMiddleware(a *service.AuthService) MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session")
			if err != nil {
				return api.NewErrorResponse(
					c,
					http.StatusUnauthorized,
					"invalid or missing session cookie",
				)
			}
			ctx := c.Request().Context()
			ok := a.IsValidSession(ctx, cookie.Value, time.Now())
			if !ok {
				a.DeleteSession(ctx, cookie.Value)
				c.SetCookie(a.DeleteSessionCookie())
				return api.NewErrorResponse(c, http.StatusUnauthorized, "invalid session cookie")
			}

			return next(c)
		}
	}
}

func AuthenticationMiddleware(svc *service.AuthService) MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session")
			if err != nil {
				return api.NewErrorResponse(
					c,
					http.StatusUnauthorized,
					"invalid authentification credentials",
				)
			}
			user, err := svc.GetCurrentUser(c.Request().Context(), cookie.Value)
			if err != nil {
				return api.NewErrorResponse(c, http.StatusUnauthorized, "invalid session cookie")
			}
			c.Set("user", *user)

			return next(c)
		}
	}
}

func AdminMiddleware() MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(domain.User)
			if !user.IsSuperuser {
				return api.NewErrorResponse(
					c,
					http.StatusForbidden,
					"Forbidden action, only available for admins",
				)
			}

			return next(c)
		}
	}
}
