package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/matiaschediek/star-trek-wheather/Core"
)

func homeLink(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "hola = %f")

}

var solarSystem = new(Core.SolarSystem)

func main() {

	today := time.Now()

	solarSystem.InitialDate = Core.Date(today.Year(), int(today.Month()), today.Day())
	solarSystem.Wheather = &Core.DaysWheather{}
	solarSystem.Ferenginar = &Core.Planet{DegreesPerDay: 1, SunDistance: 500, InitialDegrees: 90, Clockwise: true}
	solarSystem.Betazed = &Core.Planet{DegreesPerDay: 3, SunDistance: 2000, InitialDegrees: 90, Clockwise: true}
	solarSystem.Vulcano = &Core.Planet{DegreesPerDay: 5, SunDistance: 1000, InitialDegrees: 90, Clockwise: false}

	solarSystem.CalcTenYearWheather()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink).Methods("GET")
	router.HandleFunc("/wheather/day/{day}", getOneDay).Methods("GET")
	router.HandleFunc("/wheather/date/{date}", getOneDate).Methods("GET")
	router.HandleFunc("/wheather/{wheatherType}", getAllDaysByWheather).Methods("GET")
	router.HandleFunc("/wheather/allDays", getAllDays).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
func getOneDay(w http.ResponseWriter, r *http.Request) {
	day := mux.Vars(r)["day"]

	i, err := strconv.Atoi(day)
	if err != nil {
		s := fmt.Sprintf("'%s' most be the number of days after/before %s. If you need a day before use negatives numbers.", day, solarSystem.InitialDate.String())
		http.Error(w, s, http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(solarSystem.GetWheatherByDate(i))
	}

}

func getOneDate(w http.ResponseWriter, r *http.Request) {
	date := mux.Vars(r)["date"]

	layoutISO := "2006-01-02"
	t, err := time.Parse(layoutISO, date)
	if err != nil {
		http.Error(w, "The date must have the following format: 'yyyy-mm-dd'.", http.StatusBadRequest)
	} else {
		day := int(math.Round(t.Sub(solarSystem.InitialDate).Hours() / 24))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(solarSystem.GetWheatherByDate(day))
	}
}

func getAllDays(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(solarSystem.Wheather)
}

func getAllDaysByWheather(w http.ResponseWriter, r *http.Request) {

	wheatherIn := mux.Vars(r)["wheatherType"]
	s := "^[[s|S]torm|[n|N]ormal|[d|D]rought|[r|R]ainy|[o|O]ptimum]$"
	re := regexp.MustCompile(s)

	if re.MatchString(wheatherIn) {

		var daysByWheather = Core.DaysWheather{}
		all := *solarSystem.Wheather

		for _, wd := range all {
			if string(wd.Wheather) == wheatherIn {
				daysByWheather = append(daysByWheather, wd)
			}
		}
		json.NewEncoder(w).Encode(daysByWheather)

	} else {

		s := fmt.Sprintf("The value must be one of the following: %s", s)
		http.Error(w, s, http.StatusBadRequest)

	}

}
