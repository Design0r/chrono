package views

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/db/repo"
	"chrono/htmx"
	"chrono/middleware"
	"chrono/schemas"
	"chrono/service"
)

func InitProfileRoutes(group *echo.Group, db *sql.DB) {
	group.GET(
		"/profile",
		func(c echo.Context) error { return HandleProfile(c, db) },
		middleware.SessionMiddleware(db),
	)
	group.GET(
		"/profile/edit",
		func(c echo.Context) error { return HandleProfileEditForm(c, db) },
		middleware.SessionMiddleware(db),
	)

	group.PATCH(
		"/profile",
		func(c echo.Context) error { return HandleProfileEdit(c, db) },
		middleware.SessionMiddleware(db),
	)
}

func HandleProfile(c echo.Context, db *sql.DB) error {
	currUser, err := service.GetCurrentUser(db, c)
	if err != nil {
		htmx.ErrorPage(http.StatusNotFound, err.Error(), c)
		return err
	}

	notifications, err := service.GetUserNotifications(db, currUser.ID)
	if err != nil {
		htmx.ErrorPage(http.StatusNotFound, err.Error(), c)
		return err
	}
	templates.ProfilePage(currUser, notifications).Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleProfileEditForm(c echo.Context, db *sql.DB) error {
	currUser, err := service.GetCurrentUser(db, c)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	notifications, err := service.GetUserNotifications(db, currUser.ID)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}
	templates.ProfileEditForm(currUser, notifications).
		Render(context.Background(), c.Response().Writer)
	return nil
}

func HandleProfileEdit(c echo.Context, db *sql.DB) error {
	patchedData := schemas.PatchUser{}
	if err := c.Bind(&patchedData); err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	currUser, err := service.GetCurrentUser(db, c)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	updatedUser, err := service.UpdateUser(
		db,
		repo.UpdateUserParams{
			VacationDays: int64(patchedData.Vacation),
			Email:        patchedData.Email,
			Username:     patchedData.Name,
			ID:           currUser.ID,
		},
	)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	notifications, err := service.GetUserNotifications(db, currUser.ID)
	if err != nil {
		htmx.ErrorMessage(err.Error(), c)
		return err
	}

	templates.UpdateProfileWithMessage(updatedUser, notifications).
		Render(context.Background(), c.Response().Writer)
	return nil
}
