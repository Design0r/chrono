package views

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"chrono/assets/templates"
	"chrono/db/repo"
	"chrono/service"
)

func InitTokenRoutes(group *echo.Group, r *repo.Queries) {
	group.GET("/tokens", func(c echo.Context) error { return HandleTokens(c, r) })
	group.POST("/tokens", func(c echo.Context) error { return HandleCreateTokens(c, r) })
}

func HandleTokens(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)

	users, err := service.GetAllUsers(r)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	notifications, err := service.GetUserNotifications(r, currUser.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.Tokens(&currUser, notifications, users))
}

func HandleCreateTokens(c echo.Context, r *repo.Queries) error {
	currUser := c.Get("user").(repo.User)
	params, err := c.FormParams()
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	userName := params["filter"][0]
	tokenValue := params["token"][0]

	tokenNum, err := strconv.ParseFloat(tokenValue, 64)
	if err != nil {
		return RenderError(c, http.StatusBadRequest, err.Error())
	}

	user, err := service.GetUserByName(r, userName)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	_, err = service.CreateToken(r, user.ID, time.Now().Year(), tokenNum)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	msg := fmt.Sprintf("You received %v vacation token from %v", tokenNum, currUser.Username)
	_, err = service.CreateUserNotification(r, msg, user.ID)
	if err != nil {
		return RenderError(c, http.StatusInternalServerError, err.Error())
	}

	return Render(c, http.StatusOK, templates.Message("Created Token", "success"))
}
