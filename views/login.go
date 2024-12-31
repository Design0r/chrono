package views

import (
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
	return Render(c, http.StatusOK, templates.Login())
}

func HandleSignupForm(c echo.Context, r *repo.Queries) error {
	return Render(c, http.StatusOK, templates.Signup())
}

func HandleLogin(c echo.Context, r *repo.Queries) error {
	var loginUser schemas.Login
	if err := c.Bind(&loginUser); err != nil {
		return Render(c, http.StatusBadRequest, htmx.ErrorMessage("Invalid inputs", c))
	}
	user, err := service.GetUserByEmail(r, loginUser.Email)
	if err != nil {
		return Render(
			c,
			http.StatusBadRequest,
			htmx.ErrorMessage("Email or password incorrect.", c),
		)
	}

	ok := service.CheckPassword(user.Password, loginUser.Password)
	if !ok {
		return Render(
			c,
			http.StatusBadRequest,
			htmx.ErrorMessage("Email or password incorrect.", c),
		)
	}

	session, err := service.CreateSession(r, user.ID)
	if err != nil {
		return Render(
			c,
			http.StatusBadRequest,
			htmx.ErrorMessage("Internal error.", c),
		)
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
		return Render(c, http.StatusBadRequest, htmx.ErrorMessage("Invalid inputs", c))
	}

	_, err := service.GetUserByEmail(r, createUser.Email)
	if err == nil {
		return Render(
			c,
			http.StatusBadRequest,
			htmx.ErrorMessage("User with this email already exists.", c),
		)
	}

	hashedPw, err := service.HashPassword(createUser.Password)
	if err != nil {
		return Render(
			c,
			http.StatusBadRequest,
			htmx.ErrorMessage("Internal error.", c),
		)
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
		return Render(
			c,
			http.StatusBadRequest,
			htmx.ErrorMessage("Failed to create user.", c),
		)
	}

	session, err := service.CreateSession(r, user.ID)
	if err != nil {
		return Render(
			c,
			http.StatusBadRequest,
			htmx.ErrorMessage("Internal error.", c),
		)
	}

	sessionCookie := service.CreateSessionCookie(session)
	c.SetCookie(sessionCookie)

	htmx.HxRedirect("/", c)
	return nil
}
