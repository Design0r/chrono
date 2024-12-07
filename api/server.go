package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"calendar/assets"
)

type Server struct {
	router *echo.Echo
}

func NewServer(router *echo.Echo) *Server {
	return &Server{
		router: router,
	}
}

func (self *Server) InitMiddleware() {
	self.router.Use(middleware.Logger())
	// self.router.Use(
	// 	middleware.CORSWithConfig(
	// 		middleware.CORSConfig{
	// 			AllowOrigins:     []string{"http://localhost:5173"},
	// 			AllowMethods:     []string{"*"},
	// 			AllowHeaders:     []string{"*"},
	// 			AllowCredentials: true,
	// 		},
	// 	),
	// )

	staticHandler := echo.WrapHandler(http.FileServer(http.FS(assets.StaticFS)))
	self.router.GET("/static/*", staticHandler)
	self.router.Use(middleware.Recover())
}

func (self *Server) InitRoutes(group *echo.Group) {
	group.GET("", HandleIndex)
}

func (self *Server) Start(address string) error {
	return self.router.Start(address)
}
