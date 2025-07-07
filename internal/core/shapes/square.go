package shapes

import "anvil/internal/grid"

func Square(origin grid.Position, width int, height int) []grid.Position {
	size := width * height
	positions := make([]grid.Position, 0, size)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			positions = append(positions, grid.Position{X: x, Y: y})
		}
	}

	return positions
}
