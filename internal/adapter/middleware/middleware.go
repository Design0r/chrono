package middleware

import (
	"net/http"
	"time"

	"chrono/internal/adapter/htmx"
	"chrono/internal/domain"
	"chrono/internal/service"
	"chrono/schemas"

	"github.com/labstack/echo/v4"
)

type MiddlewareFunc = func(echo.HandlerFunc) echo.HandlerFunc

func SessionMiddleware(svc service.SessionService) MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session")
			if err != nil {
				return c.Redirect(http.StatusFound, "/login")
			}
			ok := svc.IsValidSession(c.Request().Context(), cookie.Value, time.Now())
			if !ok {
				svc.Delete(c.Request().Context(), cookie.Value)
				return c.Redirect(http.StatusFound, "/login")
			}

			return next(c)
		}
	}
}

func AuthenticationMiddleware(svc service.AuthService) MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session")
			if err != nil {
				return c.Redirect(http.StatusFound, "/login")
			}
			user, err := svc.GetCurrentUser(c.Request().Context(), cookie.Value)
			if err != nil {
				return c.Redirect(http.StatusFound, "/login")
			}
			c.Set("user", user)

			return next(c)
		}
	}
}

func AdminMiddleware() MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(domain.User)
			if !user.IsSuperuser {
				return htmx.RenderError(
					c,
					http.StatusForbidden,
					"Forbidden action, only available for admins",
				)
			}

			return next(c)
		}
	}
}

func HoneypotMiddleware() MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var honey schemas.Honeypot
			if err := c.Bind(&honey); err != nil {
				return htmx.RenderError(c, http.StatusBadRequest, "Invalid inputs")
			}
			if honey.IsFilled() {
				return htmx.RenderError(c, http.StatusBadRequest, "Invalid inputs")
			}

			return next(c)
		}
	}
}
