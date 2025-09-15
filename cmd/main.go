package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"

	"chrono/config"
	"chrono/db"
	"chrono/internal"
	"chrono/internal/logging"
)

func main() {
	fmt.Println(banner)
	logger, logFile, err := logging.NewTextMultiLogger("logs/chrono.log", "debug", true)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	slog.SetDefault(logger)

	slog.Info("Initializing chrono...")

	cfg := config.NewConfigFromEnv()

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.SentryUrl,
		Debug:            cfg.Debug,
		AttachStacktrace: true,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		SendDefaultPII:   true,
	}); err != nil {
		slog.Error("Sentry initialization failed", "error", err)
	}

	dbConn := db.NewDB(cfg.DbName)
	defer db.CloseDB(dbConn)

	e := echo.New()
	e.HideBanner = true

	server := internal.NewServer(e, dbConn, cfg)
	go server.Start(fmt.Sprintf(":%v", cfg.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	slog.Info("Received shutdown signal, shutting down…")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown:", "error", err)
		os.Exit(1)
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
