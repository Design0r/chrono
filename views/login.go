package views

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"calendar/assets/templates"
	"calendar/db/repo"
	"calendar/schemas"
	"calendar/service"
)

func InitLoginRoutes(group *echo.Group, db *sql.DB) {
	group.GET("/login", func(c echo.Context) error { return HandleLoginForm(c, db) })
	group.GET("/signup", func(c echo.Context) error { return HandleSignupForm(c, db) })

	group.POST("/login", func(c echo.Context) error { return HandleLogin(c, db) })
	group.POST("/signup", func(c echo.Context) error { return HandleSignup(c, db) })
	group.POST("/logout", func(c echo.Context) error { return HandleLogout(c, db) })
}

func HandleLoginForm(c echo.Context, db *sql.DB) error {
	templates.Login().Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleSignupForm(c echo.Context, db *sql.DB) error {
	templates.Signup().Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleLogin(c echo.Context, db *sql.DB) error {
	var loginUser schemas.Login
	if err := c.Bind(&loginUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid parameter")
	}
	user, err := service.GetUserByEmail(
		db,
		loginUser.Email,
	)
	if err != nil {
		return err
	}

	ok := service.CheckPassword(user.Password, loginUser.Password)
	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, "Email or Password incorrect")
	}

	session, err := service.CreateSession(db, user.ID)
	if err != nil {
		return err
	}

	sessionCookie := service.CreateSessionCookie(session)
	c.SetCookie(sessionCookie)

	templates.Home().Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleLogout(c echo.Context, db *sql.DB) error {
	cookie, err := c.Cookie("session")
	if err != nil {
		return err
	}

	sessionId, err := uuid.Parse(cookie.Value)
	if err != nil {
		return err
	}
	service.DeleteSession(db, sessionId)

	c.SetCookie(service.DeleteSessionCookie())
	return c.Redirect(http.StatusFound, "/login")
}

func HandleSignup(c echo.Context, db *sql.DB) error {
	var createUser schemas.CreateUser
	if err := c.Bind(&createUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid parameter")
	}
	hashedPw, err := service.HashPassword(createUser.Password)
	if err != nil {
		return err
	}
	user, err := service.CreateUser(
		db,
		repo.CreateUserParams{
			Username: createUser.Name,
			Email:    createUser.Email,
			Password: hashedPw,
		},
	)
	if err != nil {
		return err
	}

	session, err := service.CreateSession(db, user.ID)
	if err != nil {
		return err
	}

	sessionCookie := service.CreateSessionCookie(session)
	c.SetCookie(sessionCookie)
	return c.Redirect(http.StatusFound, "/")
}
