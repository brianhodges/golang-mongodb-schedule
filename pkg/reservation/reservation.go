package reservation
import(
    "fmt"
    "os"
    "strconv"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "golang-mongodb-schedule/pkg/location"
)

const COLLECTION = "reservations"

//easily check/throw error
func check(err error) {
    if err != nil {
        panic(err)
    }
}

type Reservation struct {
    Id bson.ObjectId `json:"id" bson:"_id,omitempty"`
    Month string
    Day int
    Year int
    Start int
    End int
    Location location.Location
    Errors  map[string]string `bson:"-"`
}

//concatenate Reservation Date
func (p Reservation) Date() string {
    return p.Month + " " + strconv.Itoa(p.Day) + ", " + strconv.Itoa(p.Year)
}

//convert Reservation (int) start to Time
func (p Reservation) Start_Time() string {
    hr := p.Start / 60
    min := p.Start % 60
    var ampm string
    if ampm = "AM"; hr >= 12 {
        ampm = "PM"
    }
    if hr > 12 {
        hr = hr - 12
    }
    if hr == 0 {
        hr = 12
    }
    return fmt.Sprintf("%02d:%02d %s", hr, min, ampm)
}

//convert Reservation (int) end to Time
func (p Reservation) End_Time() string {
    hr := p.End / 60
    min := p.End % 60
    var ampm string
    if ampm = "AM"; hr >= 12 {
        ampm = "PM"
    }
    if hr > 12 {
        hr = hr - 12
    }
    if hr == 0 {
        hr = 12
    }
    return fmt.Sprintf("%02d:%02d %s", hr, min, ampm)
}

//validate Reservation on create - booked reservations, start time < end time
func (r *Reservation) Validate() bool {
    r.Errors = make(map[string]string)
  
    if r.Start >= r.End {
      r.Errors["End"] = "End Time must be greater than Start Time"
    }
  
    session, err := mgo.Dial(os.Getenv("MONGODB_URI"))
    check(err)
    defer session.Close()
    c := session.DB(os.Getenv("MONGODB_DB")).C(COLLECTION)
    var results []Reservation
    err = c.Find(bson.M{"month": r.Month, "day": r.Day, "year": r.Year, "location": r.Location}).All(&results)

    for _, reservation := range results {
        if r.End <= reservation.Start {
            continue
        }
        if r.Start >= reservation.End {
            continue
        }
        s := fmt.Sprintf("Reservation already booked for %s on %s from %s - %s", reservation.Location.Name, reservation.Date(), reservation.Start_Time(), reservation.End_Time())
        id := fmt.Sprintf("%d", reservation.Id)
        r.Errors[id] = s
    }

    return len(r.Errors) == 0
}