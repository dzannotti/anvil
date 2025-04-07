package shapes

import "anvil/internal/grid"

func Sphere(origin grid.Position, radius int) []grid.Position {
	return Circle(origin, radius)
}

func Circle(origin grid.Position, radius int) []grid.Position {
	positions := make([]grid.Position, 0)
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			offset := grid.Position{X: x, Y: y}
			if offset.Add(origin).Distance(origin) <= radius {
				positions = append(positions, origin.Add(offset))
			}
		}
	}
	return positions
}
