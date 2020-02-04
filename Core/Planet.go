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

func (p Planet) PlanetPositionByDate(t int) Coordinates {

	var x float64
	var y float64

	InitialDegreesRad := DegreesToRadians(p.InitialDegrees)
	w := p.CalcAngularVelocity()

	rad := (w * float64(t)) + InitialDegreesRad

	x = p.SunDistance * math.Cos(rad)
	y = p.SunDistance * math.Sin(rad)

	coordinates := Coordinates{x, y}

	return coordinates
}

// Calculo de la Velocidad Angular.
// Como entrada requiere lo grados por dia que recorre el planet en cuestion.
func (p Planet) CalcAngularVelocity() float64 {
	// Calculo de dias que tarda el planeta en dar una vuelta completa.
	var t = float64(360 / p.DegreesPerDay)

	// Calculo de la velocidad angular
	w := ((float64(2) * float64(math.Pi)) / t)

	if p.Clockwise {
		w = w * (-1)
	}

	return w
}
