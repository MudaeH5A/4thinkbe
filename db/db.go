package db

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

var DB *mgo.Database

func Connection() {
	session, err := mgo.Dial(os.Getenv("MONGO_URL"))
	if err != nil {
		log.Fatal(err)
	}
	DB = session.DB(os.Getenv("MONGO_DB"))
}
