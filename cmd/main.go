package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"

	"chrono/config"
	"chrono/db"
	"chrono/service"
	"chrono/views"
)

func main() {
	fmt.Println(banner)
	log.Println("Initializing chrono...")

	cfg := config.NewConfigFromEnv()

	dbConn := db.NewDB(cfg.DbName)
	defer db.CloseDB(dbConn)

	e := echo.New()
	e.HideBanner = true

	server := views.NewServer(e, dbConn)
	server.InitMiddleware()
	server.InitRoutes()

	bot := service.NewAPIBotFromEnv()
	err := bot.Register(server.Repo)
	if err != nil {
		log.Fatal(err)
	}

	service.LoadDebugUsers(server.Repo, cfg)

	go server.Start(fmt.Sprintf(":%v", cfg.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	fmt.Println()
	log.Println("Received shutdown signal, shutting down…")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}

const banner string = `
     ________                         
    / ____/ /_  _________  ____  ____ 
   / /   / __ \/ ___/ __ \/ __ \/ __ \
  / /___/ / / / /  / /_/ / / / / /_/ /
  \____/_/ /_/_/   \____/_/ /_/\____/ 
 ======================================
 `
