package views

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"calendar/assets"
	"calendar/service"
)

type Server struct {
	Router *echo.Echo
	Db     *sql.DB
}

func NewServer(router *echo.Echo, db *sql.DB) *Server {
	return &Server{
		Router: router,
		Db:     db,
	}
}

func (self *Server) InitMiddleware() {
	self.Router.Use(middleware.Logger())

	staticHandler := echo.WrapHandler(
		http.StripPrefix(
			"/",
			http.FileServer(http.FS(assets.StaticFS)),
		),
	)
	self.Router.GET("/static/*", staticHandler)
	self.Router.Use(middleware.Recover())
	self.Router.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
}

func (self *Server) InitRoutes(group *echo.Group) {
	service.InitAPIBot(self.Db)

	InitIndexRoutes(group, self.Db)
	InitEventRoutes(group, self.Db)
	InitCalendarRoutes(group, self.Db)
	InitLoginRoutes(group, self.Db)
	InitProfileRoutes(group, self.Db)
}

func (self *Server) Start(address string) error {
	return self.Router.Start(address)
}
