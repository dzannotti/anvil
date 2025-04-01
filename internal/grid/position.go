package grid

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
