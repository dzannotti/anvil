package shapes

import "anvil/internal/grid"

func Sphere(origin grid.Position, radius int) []grid.Position {
	return Circle(origin, radius)
}
