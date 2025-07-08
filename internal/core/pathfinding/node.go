package pathfinding

import "anvil/internal/grid"

type node struct {
	pos    grid.Position
	fScore float64
}