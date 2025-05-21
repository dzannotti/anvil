package grid

import "anvil/internal/mathi"

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

// Chebyshev  Distance.
func (p Position) Distance(other Position) int {
	return mathi.Max(mathi.Abs(p.X-other.X), mathi.Abs(p.Y-other.Y))
}
