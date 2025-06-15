package internal

import (
	"database/sql"
	"log/slog"
	"time"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"chrono/config"
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
	log    *slog.Logger
	cfg    *config.Config
}

func NewServer(router *echo.Echo, db *sql.DB, cfg *config.Config) *Server {
	return &Server{
		Router: router,
		Db:     db,
		Repo:   repo.New(db),
		log:    slog.Default(),
		cfg:    cfg,
	}
}

func (s *Server) InitMiddleware() {
	s.Router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom} ${method} ${status} ${uri} ${error} ${latency_human}\n",
		CustomTimeFormat: "2006/01/02 15:04:05",
	}))

	if s.cfg.Debug == true {
		s.Router.GET("/static/*", mw.StaticHandler, mw.CacheControl)
	}

	s.Router.Use(middleware.Secure())
	s.Router.Use(middleware.Recover())
	s.Router.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	s.Router.Use(sentryecho.New(sentryecho.Options{Repanic: true}))

	s.log.Info("Initialized middleware.")
}

func (s *Server) InitRoutes() {
	userRepo := db.NewSQLUserRepo(s.Repo, s.log)
	notificationUserRepo := db.NewSQLUserNotificationRepo(s.Repo, s.log)
	notificationRepo := db.NewSQLNotificationRepo(s.Repo, s.log)
	eventRepo := db.NewSQLEventUserRepo(s.Repo, s.log)
	requestRepo := db.NewSQLRequestRepo(s.Repo, s.log)
	sessionRepo := db.NewSQLSessionRepo(s.Repo, s.log)
	refreshTokenRepo := db.NewSQLRefreshTokenRepo(s.Repo, s.log)
	vacationTokenRepo := db.NewSQLVacationTokenRepo(s.Repo, s.log)

	refreshTokenSvc := service.NewRefreshTokenService(&refreshTokenRepo, s.log)
	vacationTokenSvc := service.NewVacationTokenService(&vacationTokenRepo, s.log)
	tokenSvc := service.NewTokenService(&refreshTokenSvc, &vacationTokenSvc, s.log)
	notificationSvc := service.NewNotificationService(
		&notificationRepo,
		&notificationUserRepo,
		s.log,
	)
	userSvc := service.NewUserService(&userRepo, &notificationSvc, &tokenSvc, s.log)
	requestSvc := service.NewRequestService(&requestRepo, &userRepo, &notificationSvc, s.log)
	sessionSvc := service.NewSessionService(&sessionRepo, s.log)
	eventSvc := service.NewEventService(
		&eventRepo,
		&requestSvc,
		&userSvc,
		&vacationTokenSvc,
		s.log,
	)
	passwordHasher := auth.NewBcryptHasher(10)
	authSvc := service.NewAuthService(
		&userRepo,
		&sessionRepo,
		time.Hour*24*7,
		&passwordHasher,
		s.log,
	)

	authHandler := handler.NewAuthHandler(
		&userSvc,
		&authSvc,
		s.log,
	)
	homeHandler := handler.NewHomeHandler(&tokenSvc, &eventSvc, &notificationSvc)
	calendarHandler := handler.NewCalendarHandler(
		&userSvc,
		&notificationSvc,
		&eventSvc,
		&tokenSvc,
		s.log,
	)
	teamHandler := handler.NewTeamHandler(
		&eventSvc,
		&notificationSvc,
		&userSvc,
		s.log,
	)
	profileHandler := handler.NewProfileHandler(&userSvc, &notificationSvc, s.log)
	requestHandler := handler.NewRequestHandler(
		&requestSvc,
		&notificationSvc,
		&eventSvc,
		&vacationTokenSvc,
		s.log,
	)
	notificationHandler := handler.NewNotificationHandler(&notificationSvc, s.log)

	authGrp := s.Router.Group(
		"",
		mw.SessionMiddleware(&sessionSvc),
		mw.AuthenticationMiddleware(&authSvc),
	)
	adminGrp := authGrp.Group("", mw.AdminMiddleware())
	honeypotGrp := s.Router.Group("", mw.HoneypotMiddleware())

	calendarHandler.RegisterRoutes(authGrp)
	homeHandler.RegisterRoutes(authGrp)
	authHandler.RegisterRoutes(honeypotGrp)
	teamHandler.RegisterRoutes(authGrp, adminGrp)
	profileHandler.RegisterRoutes(authGrp, adminGrp)
	requestHandler.RegisterRoutes(adminGrp)
	notificationHandler.RegisterRoutes(authGrp)
	s.log.Info("Initialized routes.")

	bot := service.NewAPIBotFromEnv(s.log)
	bot.Register(&userSvc, &passwordHasher)
}

func (s *Server) Start(address string) error {
	s.PreStart()
	return s.Router.Start(address)
}

func (s *Server) PreStart() {
}
