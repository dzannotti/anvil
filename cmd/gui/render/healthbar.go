package ui

import (
	"anvil/internal/grid"
)

func DrawHealthbar(pos grid.Position, health, maxHealth int) {
	rect := RectFromPos(pos)
	rect.X += 10
	rect.Y -= 3
	rect.Width -= 20
	rect.Height = 6
	DrawRectangle(rect, Surface2, 2)
	rect.Width = int(float64(rect.Width) * (float64(health) / float64(maxHealth)))
	FillRectangle(rect.Expand(-1, -1), Red)
}
