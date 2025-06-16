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
	"chrono/internal/domain"
	"chrono/internal/service"
	"chrono/internal/service/auth"
)

type repos struct {
	apiCache  domain.ApiCacheRepository
	event     domain.EventRepository
	notif     domain.NotificationRepository
	notifUser domain.NotificationUserRepository
	refresh   domain.RefreshTokenRepository
	request   domain.RequestRepository
	session   domain.SessionRepository
	settings  domain.SettingsRepository
	user      domain.UserRepository
	vac       domain.VacationTokenRepository
}

type services struct {
	apiBot   service.APIBot
	auth     service.AuthService
	event    service.EventService
	holiday  service.HolidayService
	notif    service.NotificationService
	refresh  service.RefreshTokenService
	request  service.RequestService
	session  service.SessionService
	settings service.SettingsService
	token    service.TokenService
	user     service.UserService
	vac      service.VacationTokenService
	pwHasher auth.PasswordHasher
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

func (s *Server) InitRepos() {
	userRepo := db.NewSQLUserRepo(s.Repo, s.log)
	notificationUserRepo := db.NewSQLUserNotificationRepo(s.Repo, s.log)
	notificationRepo := db.NewSQLNotificationRepo(s.Repo, s.log)
	eventRepo := db.NewSQLEventUserRepo(s.Repo, s.log)
	requestRepo := db.NewSQLRequestRepo(s.Repo, s.log)
	sessionRepo := db.NewSQLSessionRepo(s.Repo, s.log)
	refreshTokenRepo := db.NewSQLRefreshTokenRepo(s.Repo, s.log)
	vacationTokenRepo := db.NewSQLVacationTokenRepo(s.Repo, s.log)

	s.repos = repos{
		user:      &userRepo,
		notifUser: &notificationUserRepo,
		notif:     &notificationRepo,
		event:     &eventRepo,
		request:   &requestRepo,
		session:   &sessionRepo,
		refresh:   &refreshTokenRepo,
		vac:       &vacationTokenRepo,
	}

	s.log.Info("Initialized repositories.")
}

func (s *Server) InitServices() {
	refreshTokenSvc := service.NewRefreshTokenService(s.repos.refresh, s.log)
	vacationTokenSvc := service.NewVacationTokenService(s.repos.vac, s.log)
	tokenSvc := service.NewTokenService(&refreshTokenSvc, &vacationTokenSvc, s.log)
	notificationSvc := service.NewNotificationService(
		s.repos.notif,
		s.repos.notifUser,
		s.log,
	)
	userSvc := service.NewUserService(s.repos.user, &notificationSvc, &tokenSvc, s.log)
	requestSvc := service.NewRequestService(s.repos.request, s.repos.user, &notificationSvc, s.log)
	sessionSvc := service.NewSessionService(s.repos.session, s.log)
	eventSvc := service.NewEventService(
		s.repos.event,
		&requestSvc,
		&userSvc,
		&vacationTokenSvc,
		s.log,
	)
	passwordHasher := auth.NewBcryptHasher(10)
	authSvc := service.NewAuthService(
		s.repos.user,
		s.repos.session,
		time.Hour*24*7,
		&passwordHasher,
		s.log,
	)

	holidaySvc := service.NewHolidayService(&userSvc, &eventSvc, s.repos.apiCache, s.log)

	s.services = services{
		refresh:  &refreshTokenSvc,
		vac:      &vacationTokenSvc,
		token:    &tokenSvc,
		notif:    &notificationSvc,
		user:     &userSvc,
		request:  &requestSvc,
		session:  &sessionSvc,
		event:    &eventSvc,
		pwHasher: &passwordHasher,
		auth:     &authSvc,
		holiday:  &holidaySvc,
	}

	s.log.Info("Initialized services.")
}

func (s *Server) InitRoutes() {
	authHandler := handler.NewAuthHandler(
		s.services.user,
		s.services.auth,
		s.log,
	)
	homeHandler := handler.NewHomeHandler(s.services.token, s.services.event, s.services.notif)
	calendarHandler := handler.NewCalendarHandler(
		s.services.user,
		s.services.notif,
		s.services.event,
		s.services.token,
		s.log,
	)
	teamHandler := handler.NewTeamHandler(
		s.services.event,
		s.services.notif,
		s.services.user,
		s.log,
	)
	profileHandler := handler.NewProfileHandler(s.services.user, s.services.notif, s.log)
	requestHandler := handler.NewRequestHandler(
		s.services.request,
		s.services.notif,
		s.services.event,
		s.services.vac,
		s.log,
	)
	notificationHandler := handler.NewNotificationHandler(s.services.notif, s.log)
	tokenHandler := handler.NewTokenHandler(
		s.services.vac,
		s.services.user,
		s.services.notif,
		s.log,
	)
	debugHandler := handler.NewDebugHandler(
		s.services.user,
		s.services.notif,
		s.services.token,
		s.services.session,
		s.log,
	)

	authGrp := s.Router.Group(
		"",
		mw.SessionMiddleware(s.services.session),
		mw.AuthenticationMiddleware(s.services.auth),
	)
	adminGrp := authGrp.Group("", mw.AdminMiddleware())
	honeypotGrp := s.Router.Group("", mw.HoneypotMiddleware())

	calendarHandler.RegisterRoutes(authGrp)
	homeHandler.RegisterRoutes(authGrp)
	notificationHandler.RegisterRoutes(authGrp)

	teamHandler.RegisterRoutes(authGrp, adminGrp)
	profileHandler.RegisterRoutes(authGrp, adminGrp)

	requestHandler.RegisterRoutes(adminGrp)
	tokenHandler.RegisterRoutes(adminGrp)
	debugHandler.RegisterRoutes(adminGrp)

	authHandler.RegisterRoutes(honeypotGrp)

	s.log.Info("Initialized routes.")
}

func (s *Server) Start(address string) error {
	s.PreStart()
	return s.Router.Start(address)
}

func (s *Server) PreStart() {
	s.InitMiddleware()
	s.InitRepos()
	s.InitServices()
	s.InitRoutes()

	bot := service.NewAPIBotFromEnv(s.log)
	bot.Register(s.services.user, s.services.pwHasher)
}
