package render

import (
	"fmt"

	"anvil/cmd/gui/ui"
	"anvil/internal/core"
	"anvil/internal/grid"
)

const cellSize = 52

func DrawWorld(w *core.World, e *core.Encounter) {
	drawGrid(w.Width(), w.Height())
	for x := 0; x < w.Width(); x++ {
		for y := 0; y < w.Height(); y++ {
			cell, _ := w.At(grid.Position{X: x, Y: y})
			drawCell(cell, e)
		}
	}
}

func drawGrid(w int, h int) {
	for y := 0; y <= h; y++ {
		ui.DrawLine(ui.Vector2i{X: 0, Y: y * cellSize}, ui.Vector2i{X: w * cellSize, Y: y * cellSize}, ui.Black, 2)
	}
	for x := 0; x <= w; x++ {
		ui.DrawLine(ui.Vector2i{X: x * cellSize, Y: 0}, ui.Vector2i{X: x * cellSize, Y: h * cellSize}, ui.Black, 2)
	}
}

func drawCell(cell *core.WorldCell, e *core.Encounter) {
	if cell.Tile == core.Wall {
		drawWall(cell.Position)
	}
	if cell.IsOccupied() {
		occupant, _ := cell.Occupant()
		drawActor(occupant, occupant == e.ActiveActor())
	}
}

func drawActor(actor *core.Actor, selected bool) {
	pos := actor.Position
	centerPos := ui.Vector2i{X: pos.X*cellSize + cellSize/2, Y: pos.Y*cellSize + cellSize/2}
	if actor.Team == core.TeamPlayers {
		ui.FillCircle(centerPos, cellSize-10, ui.Green)
	} else {
		ui.FillCircle(centerPos, cellSize-10, ui.Red)
	}
	if selected {
		ui.FillCircle(centerPos, cellSize-10, ui.Yellow)
	}
	ui.FillCircle(centerPos, cellSize-14, ui.RoyalBlue)
	shortName := fmt.Sprintf("%c%c", actor.Name[0], actor.Name[len(actor.Name)-1])
	ui.DrawString(shortName, ui.Rectangle{X: pos.X * cellSize, Y: pos.Y * cellSize, Width: cellSize, Height: cellSize}, ui.White, 15, ui.AlignMiddle)
}

func drawWall(pos grid.Position) {
	rect := ui.Rectangle{
		X:      pos.X*cellSize + 1,
		Y:      pos.Y*cellSize + 1,
		Width:  cellSize - 2,
		Height: cellSize - 2,
	}
	ui.FillRectangle(rect, ui.Brown)
}

func ToWorldPosition(pos grid.Position) ui.Vector2i {
	return ui.Vector2i{X: pos.X * cellSize, Y: pos.Y * cellSize}
}

func ToWorldPositionCenter(pos grid.Position) ui.Vector2i {
	return ui.Vector2i{X: pos.X*cellSize + cellSize/2, Y: pos.Y*cellSize + cellSize/2}
}
