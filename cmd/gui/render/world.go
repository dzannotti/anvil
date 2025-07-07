package ui

import (
	"fmt"

	"anvil/internal/core"
	"anvil/internal/grid"
)

const CellSize = 52

func DrawWorld(w *core.World, e *core.Encounter) {
	drawGrid(w.Width(), w.Height())
	for x := range w.Width() {
		for y := range w.Height() {
			pos := grid.Position{X: x, Y: y}
			FillRectangle(RectFromPos(pos).Expand(-1, -1), Surface2)
			cell := w.At(pos)
			if cell == nil {
				continue
			}
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
		occupant := cell.Occupant()
		drawActor(occupant, occupant == e.ActiveActor())
	}
	pos := fmt.Sprintf("%d,%d", cell.Position.X, cell.Position.Y)
	DrawString(pos, RectFromPos(cell.Position).Expand(-5, -5), Subtext1, 13, AlignTopLeft)
}

func drawActor(actor *core.Actor, selected bool) {
	pos := actor.Position
	centerPos := Vector2i{X: pos.X*CellSize + CellSize/2, Y: pos.Y*CellSize + CellSize/2}
	if actor.Team == core.TeamPlayers {
		FillCircle(centerPos, CellSize-10, Green)
	} else {
		FillCircle(centerPos, CellSize-10, Red)
	}
	if selected {
		FillCircle(centerPos, CellSize-10, Yellow)
	}
	FillCircle(centerPos, CellSize-14, Blue)
	shortName := fmt.Sprintf("%c%c", actor.Name[0], actor.Name[len(actor.Name)-1])
	DrawString(
		shortName,
		Rectangle{X: pos.X * CellSize, Y: pos.Y * CellSize, Width: CellSize, Height: CellSize},
		Crust,
		15,
		AlignMiddle,
	)
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
