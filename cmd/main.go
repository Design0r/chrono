package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"

	"calendar/db"
	"calendar/views"
)

func main() {
	db := db.NewDB("chrono.db")
	defer db.Close()

	e := echo.New()
	apiV1 := e.Group("")

	server := views.NewServer(e, db)
	server.InitMiddleware()
	server.InitRoutes(apiV1)

	log.Fatal(server.Start(fmt.Sprintf(":%v", os.Getenv("CHRONO_PORT"))))
}
