package definition

import "anvil/internal/grid"

type PathNode interface {
	Position() grid.Position
	SetWalkable(walkable bool)
	IsWalkable() bool
}

type Pathfinding interface {
	At(position grid.Position) (PathNode, bool)
	FindPath(start grid.Position, end grid.Position) []grid.Position
}
