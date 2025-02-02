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
		return RenderError(c, http.StatusBadRequest, "Invalid inputs")
	}
	user, err := service.GetUserByEmail(r, loginUser.Email)
	if err != nil {
		return RenderError(c, http.StatusNotFound, "Incorrect email or password")
	}

	ok := service.CheckPassword(user.Password, loginUser.Password)
	if !ok {
		return RenderError(c, http.StatusNotFound, "Incorrect email or password")
	}

	session, err := service.CreateSession(r, user.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	sessionCookie := service.CreateSessionCookie(session)
	c.SetCookie(sessionCookie)

	htmx.HxRedirect("/", c)
	return nil
}

func HandleLogout(c echo.Context, r *repo.Queries) error {
	session, err := c.Cookie("session")
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	service.DeleteSession(r, session.Value)

	c.SetCookie(service.DeleteSessionCookie())
	return c.Redirect(http.StatusFound, "/login")
}

func HandleSignup(c echo.Context, r *repo.Queries) error {
	var createUser schemas.CreateUser
	if err := c.Bind(&createUser); err != nil {
		return RenderError(c, http.StatusBadRequest, "Invalid inputs")
	}

	_, err := service.GetUserByEmail(r, createUser.Email)
	if err == nil {
		return RenderError(c, http.StatusBadRequest, "A user with this email already exists")
	}

	hashedPw, err := service.HashPassword(createUser.Password)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}
	user, err := service.CreateUser(
		r,
		repo.CreateUserParams{
			Username:     createUser.Name,
			Email:        createUser.Email,
			Color:        service.RandomHexColor(),
			VacationDays: 0,
			Password:     hashedPw,
		},
	)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	session, err := service.CreateSession(r, user.ID)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	sessionCookie := service.CreateSessionCookie(session)
	c.SetCookie(sessionCookie)

	htmx.HxRedirect("/", c)
	return nil
}
