package middleware

import (
	"net/http"
	"time"

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
			if err != nil {
				return c.Redirect(http.StatusFound, "/login")
			}
			ok := service.IsValidSession(r, cookie.Value)
			if !ok {
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

func TokenRefreshMiddleware(r *repo.Queries) MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(repo.User)
			service.InitYearlyTokens(r, user, time.Now().Year())
			return next(c)
		}
	}
}
