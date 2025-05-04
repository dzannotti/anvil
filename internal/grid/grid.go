package grid

import "anvil/internal/mathi"

type Grid[cell any] struct {
	Width  int
	Height int
	cells  []cell
}

type CellCreator[cell any] = func(pos Position) cell

func New[cell any](width int, height int, creator CellCreator[cell]) *Grid[cell] {
	cells := make([]cell, width*height)
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			cells[x+y*width] = creator(Position{X: x, Y: y})
		}
	}

	return &Grid[cell]{
		Width:  width,
		Height: height,
		cells:  cells,
	}
}

func (g Grid[cell]) fromPos(pos Position) int {
	return pos.X + pos.Y*g.Width
}

func (g Grid[cell]) fromTuple(x int, y int) int {
	return x + y*g.Width
}

func (g Grid[cell]) At(pos Position) (*cell, bool) {
	if !g.IsValidPosition(pos) {
		return nil, false
	}
	return &g.cells[g.fromPos(pos)], true
}

func (g Grid[cell]) IsValidPosition(pos Position) bool {
	return pos.X >= 0 && pos.X < g.Width && pos.Y >= 0 && pos.Y < g.Height
}

func (g Grid[cell]) CellsInRange(origin Position, radius int) []*cell {
	minX := mathi.Max(0, origin.X-radius)
	minY := mathi.Max(0, origin.Y-radius)
	maxX := mathi.Min(g.Width-1, origin.X+radius)
	maxY := mathi.Min(g.Height-1, origin.Y+radius)
	size := (maxX - minX + 1) * (maxY - minY + 1)
	cells := make([]*cell, 0, size)
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			cells = append(cells, &g.cells[g.fromTuple(x, y)])
		}
	}
	return cells
}
