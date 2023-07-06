package domain

import "math"

type Coordinate struct {
	X        int
	Y        int
	distance float64
}

func NewCoordinates(x, y int) *Coordinate {
	return &Coordinate{
		X:        x,
		Y:        y,
		distance: calculateDistance(x, y),
	}
}

func calculateDistance(x, y int) float64 {
	return math.Sqrt(math.Pow(float64(x), 2) + math.Pow(float64(y), 2))
}

func (c *Coordinate) GetDistance() float64 {
	return c.distance
}

func (c *Coordinate) GetDistanceTo(coordinate Coordinate) float64 {
	var x1 int = c.X
	var y1 int = c.Y

	var x2 int = coordinate.X
	var y2 int = coordinate.Y

	return math.Sqrt(math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2))
}
