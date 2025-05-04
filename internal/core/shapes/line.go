package shapes

import (
	"anvil/internal/grid"
	"anvil/internal/mathi"
)

func Line(from grid.Position, to grid.Position) []grid.Position {
	maxSteps := mathi.Max(mathi.Abs(to.X-from.X), mathi.Abs(to.Y-from.Y)) + 1
	result := make([]grid.Position, 0, maxSteps)

	x0, y0 := from.X, from.Y
	x1, y1 := to.X, to.Y

	dx := mathi.Abs(x1 - x0)
	dy := mathi.Abs(y1 - y0)

	sx := -1
	if x0 < x1 {
		sx = 1
	}
	sy := -1
	if y0 < y1 {
		sy = 1
	}

	err := dx - dy

	for {
		result = append(result, grid.Position{X: x0, Y: y0})

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}

	return result
}
