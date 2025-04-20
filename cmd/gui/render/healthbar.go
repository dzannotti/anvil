package render

import (
	"anvil/internal/grid"
)

func DrawHealthbar(pos grid.Position, health, maxHealth int) {
	rect := RectFromPos(pos)
	DrawRectangle(rect, Black, 2)
	rect.Width = int(float64(CellSize) * (float64(health) / float64(maxHealth)))
	FillRectangle(rect.Expand(-1, -1), Red)
}
