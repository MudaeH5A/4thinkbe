package api

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/arxdsilva/4thinkbe/db"
	"github.com/arxdsilva/4thinkbe/models"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
)

type Server struct {
	Storage *mgo.Database
	Port    string
}

func New() *Server {
	return &Server{Port: port(), Storage: db.Connection()}
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

// HomeHandler populates a new user if it does not exists
// or just returns it if It is existant in the DB
// HTTP responses:
// 200 ok
// 500 internal server error
func (s *Server) HomeHandler(c echo.Context) (err error) {
	userNumber := c.Param("userNumber")
	number, err := strconv.Atoi(userNumber)
	if err != nil {
		return
	}
	p, err := models.GetByID(s.Storage, number)
	if (err != nil) && (err != mgo.ErrNotFound) {
		return
	}
	if err == mgo.ErrNotFound {
		r := models.Room{Name: "Sala", Boxes: []models.Item{models.Item{Quantity: 2, Type: "moveis"}}}
		ca := models.Address{Street: "Santa Luiza", Number: 259, Latitude: -22.9163398, Longitude: -43.2341546}
		na := models.Address{Street: "Av Paulista", Number: 2537, Latitude: -23.5604276, Longitude: -46.6579269}
		offer := models.Offer{Distance: 50}
		p = models.Profile{ID: number, Inventory: r, NewAddress: na, CurrentAddress: ca, MovingData: time.Now(), MovingTime: time.Now(), Offer: offer}
		errC := models.Create(s.Storage, p)
		if errC != nil {
			return errC
		}
	}
	return c.JSON(http.StatusFound, p)
}
