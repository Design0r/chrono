package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/db/repo"
	"chrono/htmx"
	"chrono/service"
)

type MiddlewareFunc = func(echo.HandlerFunc) echo.HandlerFunc

func SessionMiddleware(r *repo.Queries) MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session")
			csrfCookie := service.DeleteCSRFCookie()
			if err != nil {
				c.SetCookie(csrfCookie)
				return c.Redirect(http.StatusFound, "/login")
			}
			ok := service.IsValidSession(r, cookie.Value)
			if !ok {
				c.SetCookie(csrfCookie)
				return c.Redirect(http.StatusFound, "/login")
			}

			return next(c)
		}
	}
}

func AuthenticationMiddleware(r *repo.Queries) MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := service.GetCurrentUser(r, c)
			if err != nil {
				return c.Redirect(http.StatusFound, "/login")
			}
			c.Set("user", user)

			return next(c)
		}
	}
}

func AdminMiddleware(r *repo.Queries) MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(repo.User)
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
