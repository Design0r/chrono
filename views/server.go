package views

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"chrono/assets"
	"chrono/db/repo"
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
	self.Router.Use(middleware.Logger())
	self.Router.Use(middleware.Secure())
	// self.Router.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
	// 	Skipper: func(c echo.Context) bool {
	// 		return c.Path() == "/static/*"
	// 	},
	// 	TokenLength:    32,
	// 	TokenLookup:    fmt.Sprintf("form:csrf,header:%v", echo.HeaderXCSRFToken),
	// 	ContextKey:     "csrf",
	// 	CookieName:     "_csrf",
	// 	CookieMaxAge:   86400,
	// 	CookieHTTPOnly: true,
	// 	CookieSecure:   false,
	// 	CookieDomain:   "/*",
	// }))

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

func (self *Server) InitRoutes() {
	sessionMW := self.Router.Group("", mw.SessionMiddleware(self.Repo))
	InitIndexRoutes(sessionMW, self.Repo)
	InitEventRoutes(sessionMW, self.Repo)
	InitCalendarRoutes(sessionMW, self.Repo)
	InitProfileRoutes(sessionMW, self.Repo)
	InitNotificationRoutes(sessionMW, self.Repo)
	InitRequestRoutes(sessionMW, self.Repo)
	InitTeamRoutes(sessionMW, self.Repo)

	InitLoginRoutes(self.Router.Group(""), self.Repo)
}

func (self *Server) Start(address string) error {
	return self.Router.Start(address)
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	// csrf := ctx.Get("csrf").(string)
	csrf := "1234354"
	reqCtx := context.WithValue(ctx.Request().Context(), "csrf", csrf)

	if err := t.Render(reqCtx, buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
