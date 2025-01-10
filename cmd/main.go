package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"

	"chrono/db"
	"chrono/service"
	"chrono/views"
)

func main() {
	db := db.NewDB("db/chrono.db")
	defer db.Close()

	e := echo.New()

	server := views.NewServer(e, db)
	server.InitMiddleware()
	server.InitRoutes()

	bot := service.NewAPIBotFromEnv()
	err := bot.Register(server.Repo)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(server.Start(fmt.Sprintf(":%v", os.Getenv("CHRONO_PORT"))))
}
