package util
import(
    "strconv"
    "strings"
)

type Application struct {
    Name string
    Version string
}

//easily check/throw error
func check(err error) {
    if err != nil {
        panic(err)
    }
}

//convert string Time to integer
func TimeToInteger(input string) int {
    s := strings.Split(input, ":")
    hour, tail := s[0], s[1]
    t := strings.Split(tail, " ")
    minutes, ampm := t[0], t[1]
    hr, err := strconv.Atoi(hour)
    check(err)
    min, err := strconv.Atoi(minutes)
    check(err)
    if ampm == "AM" && hr == 12 {
        hr = 0
    } else {
        if ampm == "PM" && hr != 12 {
            hr = hr + 12
        }
        hr = hr * 60 
    }
    return hr + min
}