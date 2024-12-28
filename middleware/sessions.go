package middleware

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/service"
)

type MiddlewareFunc = func(echo.HandlerFunc) echo.HandlerFunc

func SessionMiddleware(db *sql.DB) MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("session")
			if err != nil {
				return c.Redirect(http.StatusFound, "/login")
			}
			ok := service.IsValidSession(db, cookie.Value)
			if !ok {
				return c.Redirect(http.StatusFound, "/login")
			}

			return next(c)
		}
	}
}
