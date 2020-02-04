package Core

import (
	"math"
	"sort"
	"time"
)

const (
	Normal  string = "Normal"
	Drought string = "Drought"
	Rainy   string = "Rainy"
	Storm   string = "Storm"
	Optimum string = "Optimum"
)

type SolarSystem struct {
	Vulcano      *Planet
	Ferenginar   *Planet
	Betazed      *Planet
	InitialDate  time.Time
	Wheather     *DaysWheather
	maxPerimeter float64
}

type DaysWheather []DayWheather

type DayWheather struct {
	Wheather   string      `json:"Wheather"`
	Day        int         `json:"day"`
	Vulcano    Coordinates `json:"-"`
	Ferenginar Coordinates `json:"-"`
	Betazed    Coordinates `json:"-"`
	Perimeter  float64     `json:"-"`
}

func (s *SolarSystem) CalcWheatherByDate(day int) DayWheather {

	var dayWheather DayWheather = DayWheather{}

	dayWheather.Day = day

	sun := Coordinates{0, 0}

	//Pripero obtengo las coordenadas de los planetas

	dayWheather.Vulcano = s.Vulcano.PlanetPositionByDate(dayWheather.Day)
	dayWheather.Ferenginar = s.Ferenginar.PlanetPositionByDate(dayWheather.Day)
	dayWheather.Betazed = s.Betazed.PlanetPositionByDate(dayWheather.Day)

	//Calculo el área del triangulo
	area := CalcTriangleArea(dayWheather.Vulcano, dayWheather.Ferenginar, dayWheather.Betazed)

	if area == 0 {
		// Si el area es 0 significa que los tres planetas forman una recta.

		// Para saber si el sol esta en esa recta primero tenemos que obtener
		// la ecuacion de la recta. Y reemplazo los valores del Sol por la (x,y) (0,0)
		// y si la igualdad se mantiene el sol es parte de la recta.
		if CheckThePointOnStraight(dayWheather.Vulcano, dayWheather.Ferenginar, sun) {
			dayWheather.Wheather = Drought
		} else {
			dayWheather.Wheather = Optimum
		}

	} else {
		// Dado que el area es mayor a 0 los puntos forman un triangulo
		// Vamos a verificar si el Sol esta dentro del triangulo
		// Para esto vamos a calcular el area de los triangulos que forman dos de los puntos
		// por ejemplo a y b con el Sol. Y luego la divido por el area del triangulo de los planetas.
		// Lo mismo para a,c,sol  b,c,sol.
		// Si cada divicio  da entre 0 y 1 el sol esta dentro del triangulo o en uno de sus lados
		// Por el contrario se da mayor que uno el sol esta por fuera del triangulo

		areaAux := CalcTriangleArea(dayWheather.Vulcano, dayWheather.Ferenginar, sun)

		if (areaAux/area) > 1 || (areaAux/area) < 0 {
			dayWheather.Wheather = Normal
		} else {
			areaAux = CalcTriangleArea(dayWheather.Vulcano, sun, dayWheather.Betazed)
			if (areaAux/area) > 1 || (areaAux/area) < 0 {
				dayWheather.Wheather = Normal
			} else {

				areaAux = CalcTriangleArea(sun, dayWheather.Ferenginar, dayWheather.Betazed)
				if (areaAux/area) > 1 || (areaAux/area) < 0 {
					dayWheather.Wheather = Normal
				} else {

					dayWheather.Wheather = Rainy
					dayWheather.Perimeter = CalcTrianglePerimeter(dayWheather.Vulcano, dayWheather.Ferenginar, dayWheather.Betazed)
					if dayWheather.Perimeter > s.maxPerimeter {
						s.maxPerimeter = dayWheather.Perimeter
					}
				}
			}

		}

	}

	return dayWheather

}
func (s *SolarSystem) CalcTenYearWheather() {

	rainyDays := DaysWheather{}
	otherDays := DaysWheather{}
	allDays := DaysWheather{}
	y10 := s.InitialDate.AddDate(10, 0, 0)

	days := int(math.Round(y10.Sub(s.InitialDate).Hours() / 24))

	for index := 0; index < days; index++ {

		dw := s.CalcWheatherByDate(index)
		if dw.Wheather == Rainy {
			rainyDays = append(rainyDays, dw)
		} else {
			otherDays = append(otherDays, dw)
		}
	}

	for i, w := range rainyDays {
		if w.Perimeter == s.maxPerimeter {
			rainyDays[i].Wheather = Storm
		}
	}

	allDays = append(rainyDays, otherDays...)
	sort.Slice(allDays[:], func(i, j int) bool {
		return allDays[i].Day < allDays[j].Day
	})
	s.Wheather = &allDays
}

func (s *SolarSystem) GetWheatherByDate(days int) DayWheather {
	all := *s.Wheather
	for _, w := range all {
		if w.Day == days {
			return w
		}
	}

	w := s.CalcWheatherByDate(days)

	if w.Perimeter > s.maxPerimeter {
		w.Wheather = Storm
	}

	y10 := s.InitialDate.AddDate(10, 0, 0)

	d10 := int(math.Round(y10.Sub(s.InitialDate).Hours() / 24))

	if d10 > days || d10 < 0 {
		all = append(all, w)
		sort.Slice(all[:], func(i, j int) bool {
			return all[i].Day < all[j].Day
		})
		*s.Wheather = all
	}

	return w
}
