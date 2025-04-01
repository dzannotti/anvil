package pathfinding

import (
	"anvil/internal/grid"
)

type Pathfinding struct {
	width  int
	height int
	grid   *grid.Grid[Node]
}

type Result struct {
	Path []grid.Position
	Cost int
}

func New(width int, height int) *Pathfinding {
	return &Pathfinding{
		width:  width,
		height: height,
		grid: grid.New(width, height, func(pos grid.Position) Node {
			return Node{Position: pos, Walkable: true}
		}),
	}
}

func (pf *Pathfinding) At(position grid.Position) (*Node, bool) {
	return pf.grid.At(position)
}
