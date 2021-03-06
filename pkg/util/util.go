package util

import (
	"log"
	"strconv"
	"strings"
)

//Application defines application info. Used in templates
type Application struct {
	Name    string
	Version string
}

//Check logs error
func Check(err error) {
	if err != nil {
		log.Println("Error: ", err)
	}
}

//TimeToInteger converts string Time to integer
func TimeToInteger(input string) int {
	s := strings.Split(input, ":")
	hour, tail := s[0], s[1]
	t := strings.Split(tail, " ")
	minutes, ampm := t[0], t[1]
	hr, err := strconv.Atoi(hour)
	Check(err)
	min, err := strconv.Atoi(minutes)
	Check(err)
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
