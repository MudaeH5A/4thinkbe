package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
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
		offer := models.Offer{}
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
	distance, err := calculateTotalDistance(p.CurrentAddress, p.NewAddress)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	p.Offer.Distance = distance
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
// GET /:userNumber/:roomName/:boxNumber
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
	data := p.Inventory[index].Boxes[boxInt-1]
	tmp := `<style>h1,h3 {
	color: #37474f;
	text-shadow: rgba(0, 0, 0, .12) 0 0 1px;
	margin: 10px 0;
	line-height: 1.25em;
}
h1 {
	color: #FF7E56;
}
body {
    font-family: 'Open Sans', sans-serif;
    background: #fff;
    color: #76838f;
    overflow-x: hidden;
    padding-top: 56px;
}
#toolbar {
    width: 100%;
    height: 56px;
    position: fixed;
    z-index: 99;
    top: 0;
    left: 0;
    right: 0;
    justify-content: center;
    align-items: center;
    box-shadow:  0 2px 2px 0 rgba(0, 0, 0, .08);
}
</style>
<div id="toolbar">
	<h1>Mudae</h1>
</div>
	<h3> Itens da caixa </h3>
<ul>
{{range .Items}}
	<li>{{.Type}} {{.Quantity}}x</li>
{{end}}
</ul>`
	newTemplate, err := template.New("items").Parse(tmp)
	if err != nil {
		return
	}
	resp := c.Response()
	resp.Header().Set("Content-Type", "text/html")
	return newTemplate.Execute(resp.Writer, data)
}

func findRoomAndIndex(name string, rooms []models.Room) (i int, err error) {
	for i, room := range rooms {
		if room.Name == name {
			return i, nil
		}
	}
	return i, errors.New("Room not found")
}

func calculateTotalDistance(oldAddress, newAddress models.Address) (dist float64, err error) {
	origins := fmt.Sprintf("origins=%v,%v", oldAddress.Latitude, oldAddress.Longitude)
	destinations := fmt.Sprintf("destinations=%v,%v", newAddress.Latitude, newAddress.Longitude)
	apiKey := fmt.Sprintf("key=%s", os.Getenv("MAPS_KEY"))
	googleMapsEndpoint := fmt.Sprintf("https://maps.googleapis.com/maps/api/distancematrix/json?%s&%s&%s", origins, destinations, apiKey)
	resp, err := http.Get(googleMapsEndpoint)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return dist, errors.New(fmt.Sprintf("Google responded with wrong status code: %v", resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var mapsResp models.MapsResponse
	if err = json.Unmarshal(body, &mapsResp); err != nil {
		return
	}
	dist = float64(mapsResp.Rows[0].Elements[0].Distance.Value / 1000)
	return
}
