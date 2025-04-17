package render

import (
	"anvil/cmd/gui/ui"
	"anvil/internal/core"
	"anvil/internal/grid"
)

const cellSize = 64

func RenderWorld(w *core.World) {
	for y := 0; y < w.Height(); y++ {
		for x := 0; x < w.Width(); x++ {
			cell, _ := w.At(grid.Position{X: x, Y: y})
			renderCell(cell)
		}
	}
}

func renderCell(cell *core.WorldCell) {
	pos := cell.Position
	ui.FillRectangle(ui.Rectangle{X: pos.X * cellSize, Y: pos.Y * cellSize, Width: cellSize, Height: cellSize}, ui.Black)
}
