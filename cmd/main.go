package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"calendar/api"
)

func main() {
	e := echo.New()
	apiV1 := e.Group("")

	server := api.NewServer(e)
	server.InitMiddleware()
	server.InitRoutes(apiV1)

	log.Fatal(server.Start(":8080"))
}
