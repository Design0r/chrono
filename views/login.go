package views

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/db/repo"
	"chrono/htmx"
	"chrono/schemas"
	"chrono/service"
)

func InitLoginRoutes(group *echo.Group, r *repo.Queries) {
	group.GET("/login", func(c echo.Context) error { return HandleLoginForm(c, r) })
	group.GET("/signup", func(c echo.Context) error { return HandleSignupForm(c, r) })

	group.POST("/login", func(c echo.Context) error { return HandleLogin(c, r) })
	group.POST("/signup", func(c echo.Context) error { return HandleSignup(c, r) })
	group.POST("/logout", func(c echo.Context) error { return HandleLogout(c, r) })
}

func HandleLoginForm(c echo.Context, r *repo.Queries) error {
	templates.Login().Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleSignupForm(c echo.Context, r *repo.Queries) error {
	templates.Signup().Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleLogin(c echo.Context, r *repo.Queries) error {
	var loginUser schemas.Login
	if err := c.Bind(&loginUser); err != nil {
		htmx.ErrorMessage("Invalid inputs.", c)
		return err
	}
	user, err := service.GetUserByEmail(r, loginUser.Email)
	if err != nil {
		htmx.ErrorMessage("Email or password incorrect.", c)
		return err
	}

	ok := service.CheckPassword(user.Password, loginUser.Password)
	if !ok {
		htmx.ErrorMessage("Email or password incorrect.", c)
		return nil
	}

	session, err := service.CreateSession(r, user.ID)
	if err != nil {
		htmx.ErrorMessage("Internal error.", c)
		return err
	}

	sessionCookie := service.CreateSessionCookie(session)
	c.SetCookie(sessionCookie)

	htmx.HxRedirect("/", c)
	return nil
}

func HandleLogout(c echo.Context, r *repo.Queries) error {
	session, err := c.Cookie("session")
	if err != nil {
		return c.Redirect(http.StatusFound, "/error")
	}

	service.DeleteSession(r, session.Value)

	c.SetCookie(service.DeleteSessionCookie())
	return c.Redirect(http.StatusFound, "/login")
}

func HandleSignup(c echo.Context, r *repo.Queries) error {
	var createUser schemas.CreateUser
	if err := c.Bind(&createUser); err != nil {
		htmx.ErrorMessage("Invalid inputs.", c)
		return err
	}

	_, err := service.GetUserByEmail(r, createUser.Email)
	if err == nil {
		htmx.ErrorMessage("User with email already exists.", c)
		return err
	}

	hashedPw, err := service.HashPassword(createUser.Password)
	if err != nil {
		htmx.ErrorMessage("Internal error.", c)
		return err
	}
	user, err := service.CreateUser(
		r,
		repo.CreateUserParams{
			Username:     createUser.Name,
			Email:        createUser.Email,
			VacationDays: int64(createUser.Vacation),
			Password:     hashedPw,
		},
	)
	if err != nil {
		htmx.ErrorMessage("Failed to create user.", c)
		return err
	}

	session, err := service.CreateSession(r, user.ID)
	if err != nil {
		htmx.ErrorMessage("Internal error.", c)
		return err
	}

	sessionCookie := service.CreateSessionCookie(session)
	c.SetCookie(sessionCookie)

	htmx.HxRedirect("/", c)
	return nil
}
