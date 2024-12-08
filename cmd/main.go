package main

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"

	"calendar/api"
)

func main() {
	db, err := sql.Open("sqlite3", "calendar.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	e := echo.New()
	apiV1 := e.Group("")

	server := api.NewServer(e, db)
	server.InitMiddleware()
	server.InitRoutes(apiV1)

	log.Fatal(server.Start(":8080"))
}
