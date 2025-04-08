package grid

type Grid[cell any] struct {
	Width  int
	Height int
	cells  [][]cell
}

type CellCreator[cell any] = func(pos Position) cell

func New[cell any](width int, height int, creator CellCreator[cell]) *Grid[cell] {
	cells := make([][]cell, width)
	for x := 0; x < width; x++ {
		cells[x] = make([]cell, height)
		for y := 0; y < height; y++ {
			cells[x][y] = creator(Position{X: x, Y: y})
		}
	}

	return &Grid[cell]{
		Width:  width,
		Height: height,
		cells:  cells,
	}
}

func (g Grid[cell]) At(pos Position) (*cell, bool) {
	if !g.IsValidPosition(pos) {
		return nil, false
	}
	return &g.cells[pos.X][pos.Y], true
}

func (g Grid[cell]) IsValidPosition(pos Position) bool {
	return pos.X >= 0 && pos.X < g.Width && pos.Y >= 0 && pos.Y < g.Height
}

func (g Grid[cell]) CellsInRange(origin Position, radius int) []*cell {
	size := (2*radius + 1) * (2*radius + 1)
	cells := make([]*cell, 0, size)
	minP := Position{X: max(0, origin.X-radius), Y: max(0, origin.Y-radius)}
	maxP := Position{X: min(g.Width-1, origin.X+radius), Y: min(g.Height-1, origin.Y+radius)}
	for x := minP.X; x <= maxP.X; x++ {
		for y := minP.Y; y <= maxP.Y; y++ {
			cells = append(cells, &g.cells[x][y])
		}
	}
	return cells
}
