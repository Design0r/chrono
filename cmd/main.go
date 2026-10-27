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
	"github.com/labstack/echo/v5"
	_ "github.com/mattn/go-sqlite3"

	"chrono/config"
	"chrono/db"
	"chrono/internal"
	"chrono/internal/logging"
	"chrono/internal/telemetry"
)

func main() {
	fmt.Println(banner)
	log := logging.NewTextLogger("debug", true)
	slog.SetDefault(log)

	log.Info("Initializing chrono...")

	cfg := config.NewConfigFromEnv()

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.SentryUrl,
		Debug:            cfg.Debug,
		Environment:      environment(cfg.Debug),
		AttachStacktrace: true,
		EnableTracing:    true,
		TracesSampleRate: 0.1,
		// no cookies (session), request bodies (passwords) or user info like ip addresses
		DataCollection: &sentry.DataCollection{
			UserInfo:   sentry.Set(false),
			Cookies:    &sentry.KeyValueCollectionBehavior{Mode: sentry.CollectionOff},
			HTTPBodies: []sentry.BodyType{},
			HTTPHeaders: &sentry.HeaderCollectionConfig{
				Request: &sentry.KeyValueCollectionBehavior{
					Terms: []string{"forwarded", "-ip", "remote-", "via", "-user"},
				},
			},
		},
	}); err != nil {
		log.Error("Sentry initialization failed", "error", err)
	}
	defer sentry.Flush(5 * time.Second)

	telemetry.Init()
	log = telemetry.WithSentryLogs(log)
	slog.SetDefault(log)

	dbConn := db.NewDB(cfg.DbName)
	defer db.CloseDB(dbConn)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	e := echo.NewWithConfig(echo.Config{Logger: log})
	server := internal.NewServer(e, dbConn, cfg, log)
	if err := server.Start(ctx, fmt.Sprintf(":%v", cfg.Port)); err != nil {
		log.Error("Server stopped with an error", "error", err)
		sentry.Flush(5 * time.Second)
		os.Exit(1)
	}
	log.Info("Server shut down")
}

func environment(debug bool) string {
	if debug {
		return "development"
	}
	return "production"
}

const banner string = `
     ________                         
    / ____/ /_  _________  ____  ____ 
   / /   / __ \/ ___/ __ \/ __ \/ __ \
  / /___/ / / / /  / /_/ / / / / /_/ /
  \____/_/ /_/_/   \____/_/ /_/\____/ 
 ======================================
 `
