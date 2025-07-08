package shapes

import "anvil/internal/grid"

func Square(origin grid.Position, width int, height int) []grid.Position {
	size := width * height
	positions := make([]grid.Position, 0, size)
	for y := -height / 2; y < height/2; y++ {
		for x := -width / 2; x < width/2; x++ {
			positions = append(positions, origin.Add(grid.Position{X: x, Y: y}))
		}
	}

	return positions
}
