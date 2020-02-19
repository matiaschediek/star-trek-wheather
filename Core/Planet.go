package Core

import (
	"math"
)

type Planet struct {
	DegreesPerDay  float64
	SunDistance    float64
	InitialDegrees float64
	Clockwise      bool
}

const precision = 10000000

func (p *Planet) PlanetPositionByDate(t int) Coordinates {

	var x float64
	var y float64

	InitialDegreesRad := DegreesToRadians(p.InitialDegrees)

	rad := (DegreesToRadians(p.DegreesPerDay) * float64(t)) + InitialDegreesRad

	// Redondeo a 7 decimales
	x = (math.Round((p.SunDistance*math.Cos(rad))*precision) / precision)
	y = (math.Round((p.SunDistance*math.Sin(rad))*precision) / precision)

	coordinates := Coordinates{x, y}

	return coordinates
}
