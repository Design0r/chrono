package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ApiResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewErrorResponse(c echo.Context, statusCode int, message string) error {
	r := ApiResponse{
		Message: message,
		Data:    nil,
	}

	return c.JSON(statusCode, &r)
}

func NewJsonResponse(c echo.Context, data any) error {
	r := ApiResponse{
		Message: "success",
		Data:    data,
	}

	return c.JSON(http.StatusOK, &r)
}
