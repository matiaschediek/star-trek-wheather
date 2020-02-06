package Core

import (
	"math"
	"time"
)

type Coordinates struct {
	X float64
	Y float64
}

func DegreesToRadians(degrees float64) float64 {
	// Convertir grados a radianes
	return degrees * math.Pi / 180
}

func CheckThePointOnStraight(r1, r2, p Coordinates) bool {
	// Para obtener la ecuacion de la recta formada por los puntos R1 y R2 la formula es:
	// ((Y - Yr1) / (X - Xr1)) = ((Yr2 - Yr1) / (Xr2 - Xr1))
	// Por lo que si remplazamos en X e Y los valores del punto P y la igualdad se mantiene
	// Entonces el punto P pertenece a la Recta
	e1 := (p.Y - r1.Y) / (p.X - r1.X)

	e2 := (r2.Y - r1.Y) / (r2.X - r1.X)

	return e1 == e2

}

func CalcTriangleArea(a, b, c Coordinates) float64 {

	// Calculo el area del Triangulo formado las coordenas (x,y) de los puntos a, b y c
	// Este es el calculo por determinante
	// Depende del orden en que se agreguen los puntos este valor puede ser negativo por esta
	// razon necesito el modulo de ese valor
	area := math.Abs((a.X * b.Y) + (a.Y * c.X) + (b.X * c.Y) - (b.Y * c.X) - (a.Y * b.X) - (a.X * c.Y))

	// Luego del calculo del determinante se divide por dos para obtener el area
	return (area / 2)
}
func CalcTrianglePerimeter(a, b, c Coordinates) float64 {
	// Calcular distancia entre punto √ (( x - x, )² + ( y - y, )²)
	// Luego de calcular todas las distancias se suman para obtener el perimetro

	l1 := math.Sqrt(math.Pow((a.X-b.X), 2) + math.Pow((a.Y-b.Y), 2))
	l2 := math.Sqrt(math.Pow((c.X-b.X), 2) + math.Pow((c.Y-b.Y), 2))
	l3 := math.Sqrt(math.Pow((a.X-c.X), 2) + math.Pow((a.Y-c.Y), 2))
	return l1 + l2 + l3
}

func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
