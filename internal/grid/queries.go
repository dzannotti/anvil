package grid

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
