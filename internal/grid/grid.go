package grid

type Grid[cell any] struct {
	width  int
	height int
	cells  [][]cell
}

type CellCreator[cell any] = func(pos Position) cell

func New[cell any](width int, height int, creator CellCreator[cell]) *Grid[cell] {
	cells := make([][]cell, width)
	for x := 0; x < width; x++ {
		cells[x] = make([]cell, height)
		for y := 0; y < height; y++ {
			cells[x][y] = creator(NewPosition(x, y))
		}
	}

	return &Grid[cell]{
		width:  width,
		height: height,
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
	return pos.X >= 0 && pos.X < g.width && pos.Y >= 0 && pos.Y < g.height
}

func (g Grid[cell]) Width() int {
	return g.width
}

func (g Grid[cell]) Height() int {
	return g.height
}

func (g Grid[cell]) CellsInRange(origin Position, radius int) []*cell {
	cells := make([]*cell, 0)
	for x := -radius; x <= radius; x++ {
		for y := -radius; y <= radius; y++ {
			pos := origin.Add(NewPosition(x, y))
			cell, _ := g.At(pos)
			if cell != nil {
				cells = append(cells, cell)
			}
		}
	}
	return cells
}
