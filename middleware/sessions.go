package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/db/repo"
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
