package main
import (
	"fmt"
    "log"
    "os"
    "net/http"
    "html/template"
    "strconv"
    "golang-mongodb-schedule/pkg/reservation"
    "golang-mongodb-schedule/pkg/location"
    "golang-mongodb-schedule/pkg/util"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

const COLLECTION = "reservations"
const L_COLLECTION = "locations"

// Used for HTML Templates (ex. .App.Version = app.Application.Version)
type index_vars struct {
    Reservations []reservation.Reservation
    App util.Application
}

type new_vars struct {
    Locations []location.Location
    Errors map[string]string
}

//easily check/throw error
func check(err error) {
    if err != nil {
        panic(err)
    }
}

//create locations if they dont exist
func initializeLocations() {
    session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
    check(err)
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
        check(err)
        if cnt == 0 {
            err = c.Insert(l)
            check(err)
        }
    }
}

//load locations from MongoDB into array
func getLocations() []location.Location {
    session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
    check(err)
    defer session.Close()
    c := session.DB(os.Getenv("MONGODB_DB")).C(L_COLLECTION)
    var results []location.Location
    err = c.Find(nil).All(&results)
    check(err)
    return results
}

//index view handler
func index(w http.ResponseWriter, r *http.Request) {
    url := r.FormValue("url")
    app := util.Application{Name: "golang-mongodb-schedule", Version: "1.0.1"}

    session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
    check(err)
    defer session.Close()
    c := session.DB(os.Getenv("MONGODB_DB")).C(COLLECTION)
    var results []reservation.Reservation
    err = c.Find(nil).Sort("year","month","day","start").All(&results)
    check(err)
    data := index_vars{Reservations: results, App: app}
    
    if url == "" {            
        render(w, "templates/index.html", data)
    }
}

//new view handler
func new(w http.ResponseWriter, r *http.Request) {
    var results []location.Location = getLocations()
    data := new_vars{Errors: nil, Locations: results}
    render(w, "templates/new.html", data)
}

//add POST handler
func send(w http.ResponseWriter, r *http.Request) {
    day, err := strconv.Atoi(r.PostFormValue("day"))
    check(err)
    year, err := strconv.Atoi(r.PostFormValue("year"))
    check(err)
    reservation := &reservation.Reservation{
        Month: r.PostFormValue("month"),
        Day: day,
        Year: year,
        Start: util.TimeToInteger(r.PostFormValue("start")),
        End: util.TimeToInteger(r.PostFormValue("end")),
        Location: location.FindLocationByName(r.PostFormValue("location")),
    }
    
    if reservation.Validate() == false {
        results := getLocations()
        data := new_vars{Errors: reservation.Errors, Locations: results}
        render(w, "templates/new.html", data)
    } else {
        session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
        check(err)
        defer session.Close()
        c := session.DB(os.Getenv("MONGODB_DB")).C(COLLECTION)
        err = c.Insert(reservation)
        check(err)
        http.Redirect(w, r, "/", http.StatusSeeOther)
    }
}

func main() {
    initializeLocations()
	fmt.Println("Running local server @ http://localhost:" + os.Getenv("PORT"))
    http.HandleFunc("/", index)
    http.HandleFunc("/new", new)
    http.HandleFunc("/send", send)
    log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), nil))
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