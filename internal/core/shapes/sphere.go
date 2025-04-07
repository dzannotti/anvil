package shapes

import "anvil/internal/grid"

func Sphere(origin grid.Position, radius int) []grid.Position {
	return Circle(origin, radius)
}

func Circle(origin grid.Position, radius int) []grid.Position {
	size := (2*radius + 1) * (2*radius + 1)
	positions := make([]grid.Position, 0, size)
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
