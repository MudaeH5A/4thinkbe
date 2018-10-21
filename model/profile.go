package model

import "time"

type Profile struct {
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
	Distance    float64 `bson:"distance" json:"distancia"`
	LabourValue float64 `bson:"labour_value" json:"mao_de_obra"`
	KmValue     float64 `bson:"km_value" json:"valor_por_km"`
	TotalValue  float64 `bson:"total_value" json:"total_value"`
}
