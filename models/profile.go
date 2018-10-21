package models

import (
	"time"

	mgo "gopkg.in/mgo.v2"
)

type Profile struct {
	ID             int       `bson:"_id" json:"-"`
	Inventory      Room      `bson:"inventory" json:"inventario"`
	CurrentAddress Address   `bson:"current_address" json:"endereco_atual"`
	NewAddress     Address   `bson:"new_address" json:"endereco_novo"`
	MovingData     time.Time `bson:"moving_data" json:"data_mudanca"`
	MovingTime     time.Time `bson:"moving_time" json:"horario_mudanca"`
	Offer          Offer     `bson:"offer" json:"oferta"`
}

type Room struct {
	Name  string `bson:"name" json:"nome"`
	Boxes []Item `bson:"boxes" json:"caixas"`
}

type Item struct {
	Quantity int    `bson:"quantity" json:"quantidade"`
	Type     string `bson:"type" json:"tipo"`
}

type Address struct {
	Street    string  `bson:"street" json:"rua"`
	Number    int     `bson:"number" json:"numero"`
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Longitude float64 `bson:"longitude" json:"longitude"`
}

type Offer struct {
	VehicleType int     `bson:"vehicle" json:"-"`
	KmValue     float64 `bson:"km_value" json:"valor_por_km"`
	Distance    float64 `bson:"distance" json:"distancia"`
	LabourValue float64 `bson:"labour_value" json:"mao_de_obra"`
	TotalValue  float64 `bson:"total_value" json:"total_value"`
}

func (p *Profile) CreateOrUpdate(db *mgo.Database) (err error) {
	_, err = db.C("profiles").UpsertId(p.ID, p)
	return
}

func (p *Profile) DeleteByID(db *mgo.Database) error {
	return db.C("profiles").RemoveId(p.ID)
}

func GetByID(db *mgo.Database, id int) (p Profile, err error) {
	return p, db.C("profiles").FindId(id).One(&p)
}

func (o *Offer) CalculateTotalValue() {
	switch o.VehicleType {
	case 1:
		o.TotalValue = 250 + 2.0*o.Distance
		o.LabourValue = 250
		o.KmValue = 2.0
	case 2:
		o.TotalValue = 350 + 2.6*o.Distance
		o.LabourValue = 350
		o.KmValue = 2.6
	case 3:
		o.TotalValue = 500 + 3.0*o.Distance
		o.LabourValue = 500
		o.KmValue = 3.0
	}
}

func Create(db *mgo.Database, p Profile) (err error) {
	return db.C("profiles").Insert(&p)
}
