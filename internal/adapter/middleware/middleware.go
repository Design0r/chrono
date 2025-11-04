package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"chrono/assets"
	"chrono/internal/adapter/handler/api"
	"chrono/internal/adapter/htmx"
	"chrono/internal/domain"
	"chrono/internal/service"
)

type MiddlewareFunc = func(echo.HandlerFunc) echo.HandlerFunc

func SessionAPIMiddleware(svc service.SessionService) MiddlewareFunc {
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
			ok := svc.IsValidSession(c.Request().Context(), cookie.Value, time.Now())
			if !ok {
				svc.Delete(c.Request().Context(), cookie.Value)
				return api.NewErrorResponse(c, http.StatusUnauthorized, "invalid session cookie")
			}

			return next(c)
		}
	}
}

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

func AuthenticationAPIMiddleware(svc service.AuthService) MiddlewareFunc {
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
			c.Set("user", *user)

			return next(c)
		}
	}
}

func AdminAPIMiddleware() MiddlewareFunc {
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
			honey := domain.Honeypot{}
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

func SettingsAPIMiddleware(svc service.SettingsService) MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			settings, err := svc.GetFirst(c.Request().Context())
			if err != nil {
				return api.NewErrorResponse(c, http.StatusForbidden, "Failed to load settings")
			}

			ctx := context.WithValue(c.Request().Context(), "settings", settings)
			// Aktuellen Pfad auch in Context setzen
			ctx = context.WithValue(ctx, "currentPath", c.Request().URL.Path)
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)
			c.Set("settings", settings)
			return next(c)
		}
	}
}

func SettingsMiddleware(svc service.SettingsService) MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			settings, err := svc.GetFirst(c.Request().Context())
			if err != nil {
				return htmx.RenderError(
					c,
					http.StatusForbidden,
					"Failed to load settings.",
				)
			}

			ctx := context.WithValue(c.Request().Context(), "settings", settings)
			// Aktuellen Pfad auch in Context setzen
			ctx = context.WithValue(ctx, "currentPath", c.Request().URL.Path)
			req := c.Request().WithContext(ctx)
			c.SetRequest(req)
			c.Set("settings", settings)
			return next(c)
		}
	}
}
