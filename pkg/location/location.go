package location

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

//COLLECTION is the MongoDB Collection name
const COLLECTION = "locations"

//Location defines reservation location/setting
type Location struct {
	Id   bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name string
}

//FindLocationByName queries MongoDB for Location by Name
func FindLocationByName(name string) Location {
	var data Location
	session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
	util.Check(err)
	defer session.Close()
	c := session.DB(os.Getenv("MONGODB_DB")).C(COLLECTION)
	err = c.Find(bson.M{"name": name}).One(&data)
	util.Check(err)
	return data
}
