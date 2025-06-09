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

func (self *Server) InitMiddleware() {
	self.Router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom} ${method} ${status} ${uri} ${error} ${latency_human}\n",
		CustomTimeFormat: "2006/01/02 15:04:05",
	}))
	self.Router.Use(middleware.Secure())

	staticHandler := echo.WrapHandler(
		http.StripPrefix(
			"/",
			http.FileServer(http.FS(assets.StaticFS)),
		),
	)
	self.Router.GET("/static/*", staticHandler)
	self.Router.Use(middleware.Recover())
	self.Router.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	self.Router.Use(sentryecho.New(sentryecho.Options{}))
}

func (self *Server) InitRoutes() {

	userRepo := db.NewSQLUserRepo(self.Repo, slog.Default().WithGroup("user"))
	userService := service.NewUserService(&userRepo)

	//notificationUserRepo := db.NewSQLUserNotificationRepo(self.Repo, slog.Default().WithGroup("notification_user"))
	//notificationRepo := db.NewSQLNotificationRepo(self.Repo, slog.Default().WithGroup("notification"))
	//notificationService := service.NewNotificationService(&notificationRepo, &notificationUserRepo, slog.Default().WithGroup("notification"))

	//requestRepo := db.NewSQLRequestRepo(self.Repo, slog.Default().WithGroup("request"))
	//requestService := service.NewRequestService(&requestRepo, &userRepo, &notificationService, slog.Default().WithGroup("request"))

	sessionRepo := db.NewSQLSessionRepo(self.Repo, slog.Default().WithGroup("session"))
	//sessionService := service.NewSessionService(&sessionRepo, slog.Default().WithGroup("session"))

	/* 	refreshTokenRepo := db.NewSQLRefreshTokenRepo(self.Repo, slog.Default().WithGroup("refresh_token"))
	   	refreshTokenService := service.NewRefreshTokenService(&refreshTokenRepo, slog.Default().WithGroup("refresh_token"))

	   	vacationTokenRepo := db.NewSQLVacationTokenRepo(self.Repo, slog.Default().WithGroup("vacation_token"))
	   	vacationTokenService := service.NewVacationTokenService(&vacationTokenRepo, slog.Default().WithGroup("vacation_token")) */

	//tokenService := service.NewTokenService(&refreshTokenService, &vacationTokenService, slog.Default().WithGroup("token"))

	passwordHasher := auth.NewBcryptHasher(10)
	authService := service.NewAuthService(&userRepo, &sessionRepo, time.Hour*24*7, &passwordHasher, slog.Default().WithGroup("auth"))
	authHandler := handler.NewAuthHandler(&userService, &authService, slog.Default().WithGroup("auth"))

	/* 	RegisterHomeRoutes(protected)
	   	RegisterEventRoutes(protected)
	   	RegisterCalendarRoutes(protected)
	   	RegisterProfileRoutes(protected, self.Repo)
	   	RegisterNotificationRoutes(protected, self.Repo)
	   	RegisterTeamRoutes(protected, self.Repo)
		InitRequestRoutes(admin, self.Repo)
	   	InitTokenRoutes(admin, self.Repo)
	   	InitDebugRoutes(admin, self.Repo) */

	/* 	authGrp := self.Router.Group(
	   		"",
	   		mw.SessionMiddleware(&sessionService),
	   		mw.AuthenticationMiddleware(&authService),
	   	)
	   	adminGrp := authGrp.Group("", mw.AdminMiddleware()) */
	honeypotGrp := self.Router.Group("", mw.HoneypotMiddleware())

	handler.InitAuthRoutes(honeypotGrp, &authHandler)
}

func (self *Server) Start(address string) error {
	return self.Router.Start(address)
}
