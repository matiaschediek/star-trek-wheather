package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/matiaschediek/star-trek-wheather/Core"
)

func GetEnviromentVariableFloat(key string, defaultValue float64) float64 {

	valueS := os.Getenv(key)
	value, err := strconv.ParseFloat(valueS, 64)
	if err != nil {
		value = defaultValue
	}
	return value
}

func GetEnviromentVariableBool(key string, defaultValue bool) bool {

	valueS := os.Getenv(key)
	value, err := strconv.ParseBool(valueS)
	if err != nil {

		value = defaultValue
	}

	return value
}

func GetEnviromentVariableDate(key string, defaultValue time.Time) time.Time {

	valueS := os.Getenv(key)
	layoutISO := "2006-01-02"
	value, err := time.Parse(layoutISO, valueS)
	if err != nil {
		value = defaultValue
	}
	return value
}

var solarSystem = new(Core.SolarSystem)

type ResponseRange struct {
	Total int               `json:"Total"`
	Days  Core.DaysWheather `json:"Days"`
}

type ResponseRainy struct {
	Total                int               `json:"Total"`
	GreaterIntensity     int               `json:"GreaterIntensity"`
	GreaterIntensityDays Core.DaysWheather `json:"GreaterIntensityDays"`
	Days                 Core.DaysWheather `json:"Days"`
}

func main() {

	solarSystem.Wheather = new(Core.DaysWheather)

	today := time.Now()

	solarSystem.InitialDate = GetEnviromentVariableDate("SOLAR_SYSTEM_INITIAL_DATE", Core.Date(today.Year(), int(today.Month()), today.Day()))

	solarSystem.Ferenginar = &Core.Planet{
		DegreesPerDay:  GetEnviromentVariableFloat("FERENGINAR_DEGREES_PER_DAY", 1),
		SunDistance:    GetEnviromentVariableFloat("FERENGINAR_SUN_DISTANCE", 500),
		InitialDegrees: GetEnviromentVariableFloat("FERENGINAR_INITIAL_DEGREES", 90),
		Clockwise:      GetEnviromentVariableBool("FERENGINAR_CLOCKWISE", true)}
	solarSystem.Betazed = &Core.Planet{
		DegreesPerDay:  GetEnviromentVariableFloat("BETAZED_DEGREES_PER_DAY", 3),
		SunDistance:    GetEnviromentVariableFloat("BETAZED_SUN_DISTANCE", 2000),
		InitialDegrees: GetEnviromentVariableFloat("BETAZED_INITIAL_DEGREES", 90),
		Clockwise:      GetEnviromentVariableBool("BETAZED_CLOCKWISE", true)}
	solarSystem.Vulcano = &Core.Planet{
		DegreesPerDay:  GetEnviromentVariableFloat("VULCANO_DEGREES_PER_DAY", 5),
		SunDistance:    GetEnviromentVariableFloat("VULCANO_SUN_DISTANCE", 1000),
		InitialDegrees: GetEnviromentVariableFloat("VULCANO_INITIAL_DEGREES", 90),
		Clockwise:      GetEnviromentVariableBool("VULCANO_CLOCKWISE", false)}

	solarSystem.CalcTenYearWheather()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/wheather/day/{day}", getOneDay).Methods("GET")
	router.HandleFunc("/wheather/date/{date}", getOneDate).Methods("GET")
	router.HandleFunc("/wheather/type/{wheatherType}", getAllDaysByWheather).Methods("GET")
	router.HandleFunc("/wheather/range/{from}/{to}", getRange).Methods("GET")
	router.HandleFunc("/wheather/all", getAllDays).Methods("GET")
	http.ListenAndServe(":8080", router)

}
func getOneDay(w http.ResponseWriter, r *http.Request) {

	day := mux.Vars(r)["day"]

	i, err := strconv.Atoi(day)
	if err != nil {
		s := fmt.Sprintf("'%s' must be the number of days after/before %s. If you need a day before use negatives numbers.", day, solarSystem.InitialDate.String())
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

	responseRange := ResponseRange{Total: len(*solarSystem.Wheather), Days: *solarSystem.Wheather}
	json.NewEncoder(w).Encode(responseRange)
}

func getRange(w http.ResponseWriter, r *http.Request) {

	from, err := strconv.Atoi(mux.Vars(r)["from"])
	to, err := strconv.Atoi(mux.Vars(r)["to"])

	if err != nil {
		s := fmt.Sprintf("The values must be the number of days after/before %s. If you need a day before use negatives numbers.", solarSystem.InitialDate.String())
		http.Error(w, s, http.StatusBadRequest)
	} else {

		dayRange := Core.DaysWheather{}
		all := *solarSystem.Wheather

		var totalRange int
		// Si ingresan negativo en from
		if math.Signbit(float64(from)) {
			// Si "to" tambien es negativo
			if math.Signbit(float64(to)) {
				totalRange = int(math.Abs(float64(-to - (-from))))
			} else {
				totalRange = to - from
			}
		} else {

			totalRange = int(math.Abs(float64(to - from)))
		}

		if to < from {
			auxTo := to
			to = from
			from = auxTo
		}

		j := from
		for i := 0; i < totalRange+1; i++ {
			found := false
			for _, wd := range all {
				if wd.Day == j {
					dayRange = append(dayRange, wd)
					found = true
					break
				}
			}
			if !found {
				dayRange = append(dayRange, solarSystem.GetWheatherByDate(j))
			}
			j++

		}
		w.WriteHeader(http.StatusOK)
		responseRange := ResponseRange{Total: len(dayRange), Days: dayRange}
		json.NewEncoder(w).Encode(responseRange)
	}
}

func getAllDaysByWheather(w http.ResponseWriter, r *http.Request) {

	wheatherIn := mux.Vars(r)["wheatherType"]

	s := "^[[n|N][o|O][r|R][m|M][a|A][l|L]|[d|D][r|R][o|O][u|U][g|G][h|H][t|T]|[r|R][a|A][i|I][n|N][y|Y]|[o|O][p|P][t|T][i|I][m|M][u|U][m|M]]$"
	re := regexp.MustCompile(s)

	if re.MatchString(wheatherIn) {
		wheatherIn = strings.Title(strings.ToLower(wheatherIn))

		var daysByWheather = Core.DaysWheather{}

		all := *solarSystem.Wheather

		for _, wd := range all {
			if string(wd.Wheather) == wheatherIn {
				daysByWheather = append(daysByWheather, wd)
			}
		}

		if wheatherIn == "Rainy" {
			var greaterIntensityDays = Core.DaysWheather{}
			for _, wd := range daysByWheather {
				if wd.IsStorm {
					greaterIntensityDays = append(greaterIntensityDays, wd)
				}
			}
			responseRainy := ResponseRainy{Total: len(daysByWheather), Days: daysByWheather, GreaterIntensity: len(greaterIntensityDays), GreaterIntensityDays: greaterIntensityDays}
			json.NewEncoder(w).Encode(responseRainy)

		} else {
			responseRange := ResponseRange{Total: len(daysByWheather), Days: daysByWheather}
			json.NewEncoder(w).Encode(responseRange)
		}

	} else {

		s := fmt.Sprintf("The value must be one of the following: %s", s)
		http.Error(w, s, http.StatusBadRequest)

	}

}
