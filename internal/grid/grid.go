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
