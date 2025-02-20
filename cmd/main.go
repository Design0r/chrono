package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	_ "modernc.org/sqlite"

	"chrono/db"
	"chrono/service"
	"chrono/views"
)

func main() {
	fmt.Println(banner)
	log.Println("Initializing chrono...")

	db := db.NewDB("chrono.db")
	defer db.Close()

	e := echo.New()
	e.HideBanner = true

	server := views.NewServer(e, db)
	server.InitMiddleware()
	server.InitRoutes()

	bot := service.NewAPIBotFromEnv()
	err := bot.Register(server.Repo)
	if err != nil {
		log.Fatal(err)
	}

	service.LoadDebugUsers(server.Repo, "debug_users.json")

	log.Fatal(server.Start(fmt.Sprintf(":%v", os.Getenv("CHRONO_PORT"))))
}

const banner string = `
     ________                         
    / ____/ /_  _________  ____  ____ 
   / /   / __ \/ ___/ __ \/ __ \/ __ \
  / /___/ / / / /  / /_/ / / / / /_/ /
  \____/_/ /_/_/   \____/_/ /_/\____/ 
 ======================================
 `
