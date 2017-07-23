package main

import (
	"fmt"
	"golang-mongodb-schedule/pkg/location"
	"golang-mongodb-schedule/pkg/reservation"
	"golang-mongodb-schedule/pkg/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

const COLLECTION = "reservations"
const L_COLLECTION = "locations"

//IndexVars used for HTML Template Index View (ex. .App.Version = app.Application.Version)
type IndexVars struct {
	Reservations []reservation.Reservation
	App          util.Application
}

//IndexVars used for HTML Template New View
type NewVars struct {
	Locations []location.Location
	Errors    map[string]string
}

//create locations if they dont exist
func initializeLocations() {
	session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
	util.Check(err)
	defer session.Close()
	c := session.DB(os.Getenv("MONGODB_DB")).C(L_COLLECTION)
	locations := []location.Location{
		{
			Name: "Barbershop",
		},
		{
			Name: "Hair Salon",
		},
		{
			Name: "Tattoo Shop",
		},
	}
	for _, l := range locations {
		cnt, err := c.Find(bson.M{"name": l.Name}).Count()
		util.Check(err)
		if cnt == 0 {
			err = c.Insert(l)
			util.Check(err)
		}
	}
}

//load locations from MongoDB into array
func getLocations() []location.Location {
	session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
	util.Check(err)
	defer session.Close()
	c := session.DB(os.Getenv("MONGODB_DB")).C(L_COLLECTION)
	var results []location.Location
	err = c.Find(nil).All(&results)
	util.Check(err)
	return results
}

//index view handler
func index(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	app := util.Application{Name: "golang-mongodb-schedule", Version: "1.0.1"}

	session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
	util.Check(err)
	defer session.Close()
	c := session.DB(os.Getenv("MONGODB_DB")).C(COLLECTION)
	var results []reservation.Reservation
	err = c.Find(nil).Sort("year", "month", "day", "start").All(&results)
	util.Check(err)
	data := IndexVars{Reservations: results, App: app}

	if url == "" {
		render(w, "templates/index.html", data)
	}
}

//new view handler
func new(w http.ResponseWriter, r *http.Request) {
	var results []location.Location = getLocations()
	data := NewVars{Errors: nil, Locations: results}
	render(w, "templates/new.html", data)
}

//add POST handler
func send(w http.ResponseWriter, r *http.Request) {
	day, err := strconv.Atoi(r.PostFormValue("day"))
	util.Check(err)
	year, err := strconv.Atoi(r.PostFormValue("year"))
	util.Check(err)
	reservation := &reservation.Reservation{
		Month:    r.PostFormValue("month"),
		Day:      day,
		Year:     year,
		Start:    util.TimeToInteger(r.PostFormValue("start")),
		End:      util.TimeToInteger(r.PostFormValue("end")),
		Location: location.FindLocationByName(r.PostFormValue("location")),
	}

	if reservation.Validate() == false {
		results := getLocations()
		data := NewVars{Errors: reservation.Errors, Locations: results}
		render(w, "templates/new.html", data)
	} else {
		session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
		util.Check(err)
		defer session.Close()
		c := session.DB(os.Getenv("MONGODB_DB")).C(COLLECTION)
		err = c.Insert(reservation)
		util.Check(err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func main() {
	initializeLocations()
	fmt.Println("Running local server @ http://localhost:" + os.Getenv("PORT"))
	http.HandleFunc("/", index)
	http.HandleFunc("/new", new)
	http.HandleFunc("/send", send)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}

func render(w http.ResponseWriter, filename string, data interface{}) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
