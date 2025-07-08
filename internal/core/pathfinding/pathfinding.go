package pathfinding

import "anvil/internal/grid"

const (
	NormalCost   = 1.0
	DiagonalCost = 1.4
)

func FindPath(start, end grid.Position, width, height int, movementCost func(grid.Position) int) *Result {
	pf := &pathfinder{
		width:        width,
		height:       height,
		start:        start,
		end:          end,
		movementCost: movementCost,
		gCost:        make([]float64, width*height),
		cameFrom:     make([]*grid.Position, width*height),
		open:         newMinHeap(),
	}

	return pf.findPath()
}
