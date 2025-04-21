package ui

import (
	"fmt"

	"anvil/internal/core"
	"anvil/internal/grid"
)

const CellSize = 52

func DrawWorld(w *core.World, e *core.Encounter) {
	drawGrid(w.Width(), w.Height())
	for x := 0; x < w.Width(); x++ {
		for y := 0; y < w.Height(); y++ {
			pos := grid.Position{X: x, Y: y}
			FillRectangle(RectFromPos(pos).Expand(-1, -1), Surface2)
			cell, _ := w.At(pos)
			drawCell(cell, e)
		}
	}
}

func drawGrid(w int, h int) {
	for y := 0; y <= h; y++ {
		DrawLine(Vector2i{X: 0, Y: y * CellSize}, Vector2i{X: w * CellSize, Y: y * CellSize}, Overlay0, 2)
	}
	for x := 0; x <= w; x++ {
		DrawLine(Vector2i{X: x * CellSize, Y: 0}, Vector2i{X: x * CellSize, Y: h * CellSize}, Overlay0, 2)
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
	DrawString(shortName, Rectangle{X: pos.X * CellSize, Y: pos.Y * CellSize, Width: CellSize, Height: CellSize}, Crust, 15, AlignMiddle)
	DrawHealthbar(actor.Position, actor.HitPoints, actor.MaxHitPoints)
}

func drawWall(pos grid.Position) {
	rect := Rectangle{
		X:      pos.X*CellSize + 1,
		Y:      pos.Y*CellSize + 1,
		Width:  CellSize - 2,
		Height: CellSize - 2,
	}
	FillRectangle(rect, Base)
}

func ToWorldPosition(pos grid.Position) Vector2i {
	return Vector2i{X: pos.X * CellSize, Y: pos.Y * CellSize}
}

func ToWorldPositionCenter(pos grid.Position) Vector2i {
	return Vector2i{X: pos.X*CellSize + CellSize/2, Y: pos.Y*CellSize + CellSize/2}
}

func RectFromPos(pos grid.Position) Rectangle {
	return Rectangle{X: pos.X * CellSize, Y: pos.Y * CellSize, Width: CellSize, Height: CellSize}
}
