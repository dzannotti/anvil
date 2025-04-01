package pathfinding

import (
	"anvil/internal/grid"
)

type Pathfinding struct {
	width  int
	height int
	grid   *grid.Grid[Node]
}

func New(width int, height int) *Pathfinding {
	return &Pathfinding{
		width:  width,
		height: height,
		grid:   grid.New(width, height, NewNode),
	}
}

func (pf *Pathfinding) At(position grid.Position) (*Node, bool) {
	return pf.grid.At(position)
}
