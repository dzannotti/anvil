package render

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
			cell, _ := w.At(grid.Position{X: x, Y: y})
			drawCell(cell, e)
		}
	}
}

func drawGrid(w int, h int) {
	for y := 0; y <= h; y++ {
		DrawLine(Vector2i{X: 0, Y: y * CellSize}, Vector2i{X: w * CellSize, Y: y * CellSize}, Black, 2)
	}
	for x := 0; x <= w; x++ {
		DrawLine(Vector2i{X: x * CellSize, Y: 0}, Vector2i{X: x * CellSize, Y: h * CellSize}, Black, 2)
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
	if actor.Team == core.TeamPlayers {
		FillCircle(Vector2i{X: pos.X*CellSize + CellSize/2, Y: pos.Y*CellSize + CellSize/2}, CellSize-10, Green)
	} else {
		FillCircle(Vector2i{X: pos.X*CellSize + CellSize/2, Y: pos.Y*CellSize + CellSize/2}, CellSize-10, Red)
	}
	if selected {
		FillCircle(Vector2i{X: pos.X*CellSize + CellSize/2, Y: pos.Y*CellSize + CellSize/2}, CellSize-10, Yellow)
	}
	FillCircle(Vector2i{X: pos.X*CellSize + CellSize/2, Y: pos.Y*CellSize + CellSize/2}, CellSize-14, RoyalBlue)
	shortName := fmt.Sprintf("%c%c", actor.Name[0], actor.Name[len(actor.Name)-1])
	DrawString(shortName, Rectangle{X: pos.X * CellSize, Y: pos.Y * CellSize, Width: CellSize, Height: CellSize}, White, 15, AlignMiddle)
}

func drawWall(pos grid.Position) {
	rect := Rectangle{
		X:      pos.X*CellSize + 1,
		Y:      pos.Y*CellSize + 1,
		Width:  CellSize - 2,
		Height: CellSize - 2,
	}
	FillRectangle(rect, Brown)
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
