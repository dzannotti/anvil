package shapes

import (
	"math"

	"anvil/internal/grid"
)

func Circle(origin grid.Position, radius int) []grid.Position {
	capacity := int(math.Pi*float64(radius*radius)) + radius
	positions := make([]grid.Position, 0, capacity)
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				positions = append(positions, grid.Position{X: origin.X + x, Y: origin.Y + y})
			}
		}
	}

	return positions
}
