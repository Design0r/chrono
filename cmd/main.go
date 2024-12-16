package main

import (
	"log"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"

	"calendar/db"
	"calendar/views"
)

func main() {
	db := db.NewDB("calendar.db")
	defer db.Close()

	e := echo.New()
	apiV1 := e.Group("")

	server := views.NewServer(e, db)
	server.InitMiddleware()
	server.InitRoutes(apiV1)

	log.Fatal(server.Start(":8080"))
}
