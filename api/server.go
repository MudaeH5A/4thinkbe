package api

import (
	"os"

	"github.com/labstack/echo"
)

type Server struct {
	// Storage storage.Storage
	Port string
}

func New() *Server {
	return &Server{Port: port()}
}

func (s *Server) Listen() {
	e := echo.New()
	e.Static("/static", "assets")
	e.GET("/:userNumber", s.HomeHandler)
	e.Logger.Fatal(e.Start(":" + port()))
}

func port() (p string) {
	if p = os.Getenv("PORT"); p != "" {
		return
	}
	return "8080"
}

func (s *Server) HomeHandler(c echo.Context) (err error) {

	return
}
