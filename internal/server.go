package internal

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"chrono/assets"
	"chrono/db/repo"
	"chrono/internal/adapter/db"
	"chrono/internal/adapter/handler"
	mw "chrono/internal/adapter/middleware"
	"chrono/internal/service"
	"chrono/internal/service/auth"
)

type Server struct {
	Router *echo.Echo
	Db     *sql.DB
	Repo   *repo.Queries
}

func NewServer(router *echo.Echo, db *sql.DB) *Server {
	return &Server{
		Router: router,
		Db:     db,
		Repo:   repo.New(db),
	}
}

func (s *Server) InitMiddleware() {
	s.Router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom} ${method} ${status} ${uri} ${error} ${latency_human}\n",
		CustomTimeFormat: "2006/01/02 15:04:05",
	}))
	s.Router.Use(middleware.Secure())

	staticHandler := echo.WrapHandler(
		http.StripPrefix(
			"/",
			http.FileServer(http.FS(assets.StaticFS)),
		),
	)
	s.Router.GET("/static/*", staticHandler)
	s.Router.Use(middleware.Recover())
	s.Router.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	s.Router.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
}

func (s *Server) InitRoutes() {

	userRepo := db.NewSQLUserRepo(s.Repo, slog.Default().WithGroup("user"))
	notificationUserRepo := db.NewSQLUserNotificationRepo(s.Repo, slog.Default().WithGroup("notification_user"))
	notificationRepo := db.NewSQLNotificationRepo(s.Repo, slog.Default().WithGroup("notification"))
	eventRepo := db.NewSQLEventUserRepo(s.Repo, slog.Default().WithGroup("event"))
	requestRepo := db.NewSQLRequestRepo(s.Repo, slog.Default().WithGroup("request"))
	sessionRepo := db.NewSQLSessionRepo(s.Repo, slog.Default().WithGroup("session"))
	refreshTokenRepo := db.NewSQLRefreshTokenRepo(s.Repo, slog.Default().WithGroup("refresh_token"))
	vacationTokenRepo := db.NewSQLVacationTokenRepo(s.Repo, slog.Default().WithGroup("vacation_token"))

	userService := service.NewUserService(&userRepo)
	notificationService := service.NewNotificationService(&notificationRepo, &notificationUserRepo, slog.Default().WithGroup("notification"))
	requestService := service.NewRequestService(&requestRepo, &userRepo, &notificationService, slog.Default().WithGroup("request"))
	sessionService := service.NewSessionService(&sessionRepo, slog.Default().WithGroup("session"))
	refreshTokenService := service.NewRefreshTokenService(&refreshTokenRepo, slog.Default().WithGroup("refresh_token"))
	vacationTokenService := service.NewVacationTokenService(&vacationTokenRepo, slog.Default().WithGroup("vacation_token"))
	tokenService := service.NewTokenService(&refreshTokenService, &vacationTokenService, slog.Default().WithGroup("token"))
	eventService := service.NewEventService(&eventRepo, &requestService, &userService, &vacationTokenService, slog.Default().WithGroup("event"))
	passwordHasher := auth.NewBcryptHasher(10)
	authService := service.NewAuthService(&userRepo, &sessionRepo, time.Hour*24*7, &passwordHasher, slog.Default().WithGroup("auth"))

	authHandler := handler.NewAuthHandler(&userService, &authService, slog.Default().WithGroup("auth"))
	homeHandler := handler.NewHomeHandler(&tokenService, &eventService, &notificationService)

	authGrp := s.Router.Group(
		"",
		mw.SessionMiddleware(&sessionService),
		mw.AuthenticationMiddleware(&authService),
	)
	//adminGrp := authGrp.Group("", mw.AdminMiddleware())
	honeypotGrp := s.Router.Group("", mw.HoneypotMiddleware())

	homeHandler.RegisterRoutes(authGrp)
	authHandler.RegisterRoutes(honeypotGrp)
}

func (s *Server) Start(address string) error {
	return s.Router.Start(address)
}
