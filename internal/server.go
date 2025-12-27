package internal

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"chrono/config"
	"chrono/db/repo"
	"chrono/internal/adapter/db"
	"chrono/internal/adapter/handler/api"
	mw "chrono/internal/adapter/middleware"
	"chrono/internal/domain"
	"chrono/internal/service"
	"chrono/internal/service/auth"
)

type repos struct {
	apiCache   domain.ApiCacheRepository
	event      domain.EventRepository
	notif      domain.NotificationRepository
	notifUser  domain.NotificationUserRepository
	refresh    domain.RefreshTokenRepository
	request    domain.RequestRepository
	session    domain.SessionRepository
	settings   domain.SettingsRepository
	user       domain.UserRepository
	vac        domain.VacationTokenRepository
	timestamps domain.TimestampsRepository
}

type services struct {
	apiBot     *service.APIBot
	auth       *service.AuthService
	event      *service.EventService
	holiday    *service.HolidayService
	notif      *service.NotificationService
	request    *service.RequestService
	settings   *service.SettingsService
	token      *service.TokenService
	user       *service.UserService
	pwHasher   auth.PasswordHasher
	krank      *service.KrankheitsExport
	awork      *service.AworkService
	timestamps *service.TimestampsService
}

type Server struct {
	Router   *echo.Echo
	Db       *sql.DB
	Repo     *repo.Queries
	log      *slog.Logger
	cfg      *config.Config
	repos    repos
	services services
}

func NewServer(router *echo.Echo, db *sql.DB, cfg *config.Config, log *slog.Logger) *Server {
	return &Server{
		Router: router,
		Db:     db,
		Repo:   repo.New(db),
		log:    log,
		cfg:    cfg,
	}
}

func (s *Server) InitMiddleware() {
	s.Router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom} ${method} ${status} ${uri} ${error} ${latency_human}\n",
		CustomTimeFormat: "2006/01/02 15:04:05",
	}))
	s.Router.Use(middleware.Secure())
	s.Router.Use(middleware.Recover())
	s.Router.Use(
		middleware.CORSWithConfig(
			middleware.CORSConfig{
				AllowOrigins: []string{
					"http://localhost:8080",
					"http://localhost:5173",
					"http://192.168.0.35:5173",
					"https://chrono.theapic.com",
				},
				AllowCredentials: true,
			},
		),
	)
	s.Router.Use(sentryecho.New(sentryecho.Options{Repanic: true}))

	s.log.Info("Initialized middleware.")
}

func (s *Server) InitRepos() {
	userRepo := db.NewSQLUserRepo(s.Repo, s.log)
	notificationUserRepo := db.NewSQLUserNotificationRepo(s.Repo, s.log)
	notificationRepo := db.NewSQLNotificationRepo(s.Repo, s.log)
	eventRepo := db.NewSQLEventUserRepo(s.Repo, s.log)
	requestRepo := db.NewSQLRequestRepo(s.Repo, s.log)
	sessionRepo := db.NewSQLSessionRepo(s.Repo, s.log)
	refreshTokenRepo := db.NewSQLRefreshTokenRepo(s.Repo, s.log)
	vacationTokenRepo := db.NewSQLVacationTokenRepo(s.Repo, s.log)
	apiCacheRepo := db.NewSQLAPICacheRepo(s.Repo, s.log)
	settingsRepo := db.NewSQLSettingsRepo(s.Repo, s.log)
	timestampsRepo := db.NewSQLTimestampsRepo(s.Repo, s.log)

	s.repos = repos{
		user:       userRepo,
		notifUser:  notificationUserRepo,
		notif:      notificationRepo,
		event:      eventRepo,
		request:    requestRepo,
		session:    sessionRepo,
		refresh:    refreshTokenRepo,
		vac:        vacationTokenRepo,
		apiCache:   apiCacheRepo,
		settings:   settingsRepo,
		timestamps: timestampsRepo,
	}

	s.log.Info("Initialized repositories.")
}

func (s *Server) InitServices() {
	tokenSvc := service.NewTokenService(s.repos.refresh, s.repos.vac, s.log)
	notificationSvc := service.NewNotificationService(
		s.repos.notif,
		s.repos.notifUser,
		s.log,
	)
	userSvc := service.NewUserService(s.repos.user, notificationSvc, tokenSvc, s.log)
	requestSvc := service.NewRequestService(s.repos.request, s.repos.user, notificationSvc, s.log)
	eventSvc := service.NewEventService(
		s.repos.event,
		requestSvc,
		userSvc,
		tokenSvc,
		s.log,
	)
	passwordHasher := auth.NewBcryptHasher(10)
	authSvc := service.NewAuthService(
		s.repos.user,
		s.repos.session,
		time.Hour*24*7,
		!s.cfg.Debug,
		passwordHasher,
		s.log,
	)

	holidaySvc := service.NewHolidayService(userSvc, eventSvc, s.repos.apiCache, s.log)
	settingSvc := service.NewSettingsService(s.repos.settings, s.log)
	krankSvc := service.NewKrankheitsExportService(eventSvc, userSvc)
	aworkSvc := service.NewAworkService(eventSvc, userSvc, s.log)
	timestampSvc := service.NewTimestampsService(s.repos.timestamps, eventSvc, s.log)

	s.services = services{
		token:      tokenSvc,
		notif:      notificationSvc,
		user:       userSvc,
		request:    requestSvc,
		event:      eventSvc,
		pwHasher:   passwordHasher,
		auth:       authSvc,
		holiday:    holidaySvc,
		settings:   settingSvc,
		krank:      krankSvc,
		awork:      aworkSvc,
		timestamps: timestampSvc,
	}

	s.log.Info("Initialized services.")
}

func (s *Server) InitAPIRoutes() {
	authHandler := api.NewAPIAuthHandler(
		s.services.user,
		s.services.auth,
		s.log,
	)
	userHandler := api.NewAPIUserHandler(
		s.services.user,
		s.services.event,
		s.services.auth,
		s.services.token,
		s.log,
	)
	eventHandler := api.NewAPIEventHandler(
		s.services.user,
		s.services.event,
		s.services.token,
		s.log,
	)
	requestHandler := api.NewAPIRequestsHandler(
		s.services.request,
		s.services.event,
		s.services.token,
		s.log,
	)
	tokenHandler := api.NewAPITokenHandler(
		s.services.token,
		s.services.user,
		s.services.notif,
		s.log,
	)
	settingsHandler := api.NewAPISettingsHandler(s.services.settings)
	exportHander := api.NewAPIExportHandler(s.services.krank)
	aworkHandler := api.NewAPIAworkHandler(
		s.services.user,
		s.services.event,
		s.services.awork,
		s.log,
	)
	notificationHandler := api.NewAPINotificationHandler(s.services.notif, s.log)
	timestampsHandler := api.NewAPITimestampsHandler(s.services.timestamps)

	apiGrp := s.Router.Group("/api/v1")
	authGrp := apiGrp.Group(
		"",
		mw.SessionMiddleware(s.services.auth),
		mw.AuthenticationMiddleware(s.services.auth),
	)

	adminGrp := authGrp.Group("", mw.AdminMiddleware())

	authHandler.RegisterRoutes(apiGrp)

	userHandler.RegisterRoutes(authGrp)
	eventHandler.RegisterRoutes(authGrp)
	aworkHandler.RegisterRoutes(authGrp)
	notificationHandler.RegisterRoutes(authGrp)
	timestampsHandler.RegisterRoutes(authGrp, adminGrp)

	requestHandler.RegisterRoutes(adminGrp)
	tokenHandler.RegisterRoutes(adminGrp)
	settingsHandler.RegisterRoutes(adminGrp)
	exportHander.RegisterRoutes(adminGrp)

	apiGrp.GET(
		"/health",
		func(c echo.Context) error { return c.NoContent(http.StatusOK) },
	)

	s.log.Info("Initialized api routes.")
}

func (s *Server) Start(address string) error {
	err := s.PreStart()
	if err != nil {
		s.log.Error("Failed to run PreStart.", "error", err.Error())
		return err
	}
	return s.Router.Start(address)
}

func (s *Server) PreStart() error {
	s.InitMiddleware()
	s.InitRepos()
	s.InitServices()
	s.InitAPIRoutes()

	settings := domain.Settings{SignupEnabled: false}
	_, err := s.services.settings.Init(context.Background(), settings)
	if err != nil {
		s.log.Error("Failed to init settings.")
		return err
	}

	bot := service.NewAPIBotFromEnv(s.log)
	bot.Register(context.Background(), s.services.user, s.services.pwHasher)

	return nil
}
