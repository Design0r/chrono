package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/assets"
	"chrono/db/repo"
	"chrono/htmx"
	"chrono/schemas"
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
				service.DeleteSession(r, cookie.Value)
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

func HoneypotMiddleware(r *repo.Queries) MiddlewareFunc {
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

func CacheControl(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "public, max-age=86400") // 1 day
		return next(c)
	}
}

var StaticHandler = echo.WrapHandler(
	http.StripPrefix(
		"/",
		http.FileServer(http.FS(assets.StaticFS)),
	),
)
