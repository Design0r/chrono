package views

import (
	"database/sql"
	"net/http"

	"github.com/a-h/templ"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"chrono/assets"
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

func (self *Server) InitMiddleware() {
	self.Router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom} ${method} ${status} ${uri} ${error} ${latency_human}\n",
		CustomTimeFormat: "2006/01/02 15:04:05",
	}))
	self.Router.Use(middleware.Secure())

	cacheStatic := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "public, max-age=86400") // 1 day
			return next(c)
		}
	}

	staticHandler := echo.WrapHandler(
		http.StripPrefix(
			"/",
			http.FileServer(http.FS(assets.StaticFS)),
		),
	)
	self.Router.GET("/static/*", staticHandler, cacheStatic)
	self.Router.Use(middleware.Recover())
	self.Router.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	self.Router.Use(sentryecho.New(sentryecho.Options{}))
}

func (self *Server) InitRoutes() {
	protected := self.Router.Group(
		"",
		mw.SessionMiddleware(self.Repo),
		mw.AuthenticationMiddleware(self.Repo),
	)
	InitHomeRoutes(protected, self.Repo)
	InitEventRoutes(protected, self.Repo)
	InitCalendarRoutes(protected, self.Repo)
	InitProfileRoutes(protected, self.Repo)
	InitNotificationRoutes(protected, self.Repo)
	InitTeamRoutes(protected, self.Repo)

	admin := protected.Group("", mw.AdminMiddleware(self.Repo))
	InitRequestRoutes(admin, self.Repo)
	InitTokenRoutes(admin, self.Repo)
	InitDebugRoutes(admin, self.Repo)

	honeyot := self.Router.Group("", mw.HoneypotMiddleware(self.Repo))
	InitLoginRoutes(honeyot, self.Repo)
}

func (self *Server) Start(address string) error {
	return self.Router.Start(address)
}

func Render(c echo.Context, statusCode int, t templ.Component) error {
	return htmx.Render(c, statusCode, t)
}

func RenderError(c echo.Context, statusCode int, msg string) error {
	return htmx.RenderError(c, statusCode, msg)
}
