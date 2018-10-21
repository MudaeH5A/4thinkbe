package db

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

func Connection() *mgo.Database {
	session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
	if err != nil {
		log.Fatal(err)
	}
	return session.DB("mudae")
}
