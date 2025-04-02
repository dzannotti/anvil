package shapes

import "anvil/internal/grid"

func Sphere(origin grid.Position, radius int) []grid.Position {
	positions := make([]grid.Position, 0)
	for x := -radius; x <= radius; x++ {
		for y := -radius; y <= radius; y++ {
			offset := grid.Position{X: x, Y: y}
			if offset.Distance(origin) <= radius {
				positions = append(positions, origin.Add(offset))
			}
		}
	}
	return positions
}
