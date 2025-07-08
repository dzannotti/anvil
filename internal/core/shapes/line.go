package shapes

import (
	"anvil/internal/grid"
	"anvil/internal/mathi"
)

func Line(from, to grid.Position) []grid.Position {
	deltaX := to.X - from.X
	deltaY := to.Y - from.Y

	steps := mathi.Max(mathi.Abs(deltaX), mathi.Abs(deltaY))
	if steps == 0 {
		return []grid.Position{from}
	}

	result := make([]grid.Position, 0, steps+1)

	for i := 0; i <= steps; i++ {
		currentX := from.X + (deltaX*i)/steps
		currentY := from.Y + (deltaY*i)/steps
		result = append(result, grid.Position{X: currentX, Y: currentY})
	}

	return result
}
