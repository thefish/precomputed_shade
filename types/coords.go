package types

import "math"

type Coords struct {
	X, Y int
}

func (c *Coords) Get() (int, int) {
	return c.X, c.Y
}

func (c *Coords) DistanceTo(o Coords) float64 {
	dx := c.X - o.X
	dy := c.X - o.Y
	return math.Sqrt(math.Pow(float64(dx), 2) + math.Pow(float64(dy), 2))
}

func (c *Coords) IsAdjacentTo(o *Coords) bool {
	var xDiff, yDiff int
	if c.X > o.X {
		xDiff = c.X - o.X
	} else {
		xDiff = o.X - c.X
	}
	if c.Y > o.Y {
		yDiff = c.Y - o.Y
	} else {
		yDiff = o.Y - c.Y
	}
	return xDiff < 2 && yDiff < 2
}
