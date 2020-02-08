package Core

import (
	"math"
)

type Planet struct {
	DegreesPerDay   float64
	SunDistance     float64
	InitialDegrees  float64
	Clockwise       bool
	AngularVelocity float64
}

func (p *Planet) PlanetPositionByDate(t int) Coordinates {

	var x float64
	var y float64

	InitialDegreesRad := DegreesToRadians(p.InitialDegrees)

	rad := (p.AngularVelocity * float64(t)) + InitialDegreesRad

	// Redondeo a 7 decimales
	x = (math.Round((p.SunDistance*math.Cos(rad))*10000000) / 10000000)
	y = (math.Round((p.SunDistance*math.Sin(rad))*10000000) / 10000000)

	coordinates := Coordinates{x, y}

	return coordinates
}

// Calculo de la Velocidad Angular.
// Como entrada requiere lo grados por dia que recorre el planet en cuestion.
func (p *Planet) CalcAngularVelocity() {
	// Calculo de dias que tarda el planeta en dar una vuelta completa.
	var t = float64(360 / p.DegreesPerDay)

	// Calculo de la velocidad angular
	w := ((float64(2) * float64(math.Pi)) / t)

	if p.Clockwise {
		w = w * (-1)
	}

	p.AngularVelocity = w
}
