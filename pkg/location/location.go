package location

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

const COLLECTION = "locations"

//easily check/throw error
func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Location struct {
	Id   bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name string
}

//query MongoDB for Location by Name
func FindLocationByName(name string) Location {
	var data Location
	session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
	check(err)
	defer session.Close()
	c := session.DB(os.Getenv("MONGODB_DB")).C(COLLECTION)
	err = c.Find(bson.M{"name": name}).One(&data)
	check(err)
	return data
}
