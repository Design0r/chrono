package views

import (
	"database/sql"

	"github.com/a-h/templ"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"chrono/db/repo"
	"chrono/htmx"
	mw "chrono/middleware"
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
	s.Router.Use(sentryecho.New(sentryecho.Options{Repanic: true}))
	s.Router.GET("/static/*", mw.StaticHandler, mw.CacheControl)
	s.Router.Use(middleware.Recover())
	s.Router.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
}

func (s *Server) InitRoutes() {
	protected := s.Router.Group(
		"",
		mw.SessionMiddleware(s.Repo),
		mw.AuthenticationMiddleware(s.Repo),
	)
	InitHomeRoutes(protected, s.Repo)
	InitEventRoutes(protected, s.Repo)
	InitCalendarRoutes(protected, s.Repo)
	InitProfileRoutes(protected, s.Repo)
	InitNotificationRoutes(protected, s.Repo)
	InitTeamRoutes(protected, s.Repo)

	admin := protected.Group("", mw.AdminMiddleware(s.Repo))
	InitRequestRoutes(admin, s.Repo)
	InitTokenRoutes(admin, s.Repo)
	InitDebugRoutes(admin, s.Repo)

	honeyot := s.Router.Group("", mw.HoneypotMiddleware(s.Repo))
	InitLoginRoutes(honeyot, s.Repo)
}

func (s *Server) Start(address string) error {
	return s.Router.Start(address)
}

func Render(c echo.Context, statusCode int, t templ.Component) error {
	return htmx.Render(c, statusCode, t)
}

func RenderError(c echo.Context, statusCode int, msg string) error {
	return htmx.RenderError(c, statusCode, msg)
}
