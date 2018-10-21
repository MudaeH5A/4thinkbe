package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/arxdsilva/4thinkbe/db"
	"github.com/arxdsilva/4thinkbe/models"
	"github.com/labstack/echo"
	qrcode "github.com/skip2/go-qrcode"
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
	e.GET("/:userNumber/:room/:boxNumber/code", s.BoxCoder)
	e.GET("/:userNumber/:room/:boxNumber", s.BoxContent)
	e.POST("/:userNumber/:vehicle", s.VehicleHandler)
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
// GET /:userNumber
// HTTP responses:
// 200 ok
// 400 bad request
// 500 internal server error
func (s *Server) HomeHandler(c echo.Context) (err error) {
	userNumber := c.Param("userNumber")
	number, err := strconv.Atoi(userNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	p, err := models.GetByID(s.Storage, number)
	if (err != nil) && (err != mgo.ErrNotFound) {
		return echo.NewHTTPError(500, err)
	}
	if err == mgo.ErrNotFound {
		r := models.Room{
			Name: "sala",
			Boxes: []models.Box{
				models.Box{
					Items: []models.Item{
						models.Item{
							Quantity: 2,
							Type:     "moveis",
						},
					},
				},
				models.Box{
					Items: []models.Item{
						models.Item{
							Quantity: 1,
							Type:     "tv",
						},
					},
				},
			},
		}
		ca := models.Address{Street: "Santa Luiza", Number: 259, Latitude: -22.9163398, Longitude: -43.2341546}
		na := models.Address{Street: "Av Paulista", Number: 2537, Latitude: -23.5604276, Longitude: -46.6579269}
		offer := models.Offer{Distance: 50}
		p = models.Profile{
			ID:             number,
			Inventory:      []models.Room{r},
			NewAddress:     na,
			CurrentAddress: ca,
			MovingData:     time.Now(),
			MovingTime:     time.Now(),
			Offer:          offer,
		}
		errC := models.Create(s.Storage, p)
		if errC != nil {
			return echo.NewHTTPError(500, errC)
		}
	}
	return c.JSON(http.StatusFound, p)
}

// VehicleHandler handles the user offer infos
// it helps adding data to the user
// POST /:userNumber/:vehicleNumber [1-3]
//
// HTTP responses:
// 201 created
// 400 bad request
// 500 internal server error
func (s *Server) VehicleHandler(c echo.Context) (err error) {
	userNumber := c.Param("userNumber")
	number, err := strconv.Atoi(userNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	vehicleNumber := c.Param("vehicle")
	vehicle, err := strconv.Atoi(vehicleNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if (vehicle > 3) || (vehicle < 1) {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("Vehicle type can only be between [1-3]"))
	}
	p, err := models.GetByID(s.Storage, number)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	p.Offer.VehicleType = vehicle
	p.Offer.CalculateTotalValue()
	err = p.CreateOrUpdate(s.Storage)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, p.Offer)
}

// BoxCoder generates QR Codes for a specific box
// GET /:userNumber/:room/:box/code
//
// HTTP responses:
// 200 OK
// 500 internal server error
func (s *Server) BoxCoder(c echo.Context) (err error) {
	userNumber := c.Param("userNumber")
	room := c.Param("room")
	boxNumber := c.Param("boxNumber")
	appURL := fmt.Sprintf("https://mudae.herokuapp.com/%s/%s/%s", userNumber, room, boxNumber)
	png, err := qrcode.Encode(appURL, qrcode.Medium, 256)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	resp := c.Response()
	resp.Header().Set("Content-Type", "image/png")
	resp.Header().Set("Content-Length", strconv.Itoa(len(png)))
	_, err = resp.Write(png)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return
}

// BoxContent lists the contents of a specific
// box after the QR code call to the API
// GET /:userNumber/:room/:box
//
// HTTP responses:
// 302 OK
// 400 bad request
// 404 not found
func (s *Server) BoxContent(c echo.Context) (err error) {
	userNumber := c.Param("userNumber")
	room := c.Param("room")
	boxNumber := c.Param("boxNumber")
	number, err := strconv.Atoi(userNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	p, err := models.GetByID(s.Storage, number)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err)
	}
	index, err := findRoomAndIndex(room, p.Inventory)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	boxInt, err := strconv.Atoi(boxNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusFound, p.Inventory[index].Boxes[boxInt])
}

func findRoomAndIndex(name string, rooms []models.Room) (i int, err error) {
	for i, room := range rooms {
		if room.Name == name {
			return i, nil
		}
	}
	return i, errors.New("Room not found")
}
