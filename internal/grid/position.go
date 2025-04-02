package grid

import "github.com/adam-lavrik/go-imath/ix"

type Position struct {
	X int
	Y int
}

func (p Position) Add(other Position) Position {
	return Position{
		X: p.X + other.X,
		Y: p.Y + other.Y,
	}
}

func (p Position) Subtract(other Position) Position {
	return Position{
		X: p.X - other.X,
		Y: p.Y - other.Y,
	}
}

// Chebyshev  Distance
func (p Position) Distance(other Position) int {
	return max(ix.Abs(p.X-other.X), ix.Abs(p.Y-other.Y))
}
