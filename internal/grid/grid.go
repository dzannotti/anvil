package grid

type Grid[cell any] struct {
	Width  int
	Height int
	cells  []cell
}

type CellCreator[cell any] = func(pos Position) cell

func New[cell any](width int, height int, creator CellCreator[cell]) *Grid[cell] {
	cells := make([]cell, width*height)
	for x := range width {
		for y := range height {
			cells[x+y*width] = creator(Position{X: x, Y: y})
		}
	}

	return &Grid[cell]{
		Width:  width,
		Height: height,
		cells:  cells,
	}
}

func (g *Grid[cell]) fromXY(x int, y int) int {
	return x + y*g.Width
}

func (g *Grid[cell]) At(pos Position) *cell {
	return &g.cells[g.fromXY(pos.X, pos.Y)]
}

func (g *Grid[cell]) IsValidPosition(pos Position) bool {
	return pos.X >= 0 && pos.X < g.Width && pos.Y >= 0 && pos.Y < g.Height
}

func (g *Grid[cell]) Cells(positions []Position) []*cell {
	cells := make([]*cell, 0, len(positions))
	for _, pos := range positions {
		if !g.IsValidPosition(pos) {
			continue
		}

		cells = append(cells, g.At(pos))
	}

	return cells
}
