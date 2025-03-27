package definition

import "anvil/internal/grid"

type WorldCell interface {
	Position() grid.Position
	Occupant() (Creature, bool)
	AddOccupant(Creature)
	RemoveOccupant(Creature)
	IsOccupied() bool
}

type World interface {
	At(grid.Position) (WorldCell, bool)
	Navigation() Pathfinding
	CreaturesInRange(grid.Position, int) []Creature
	Width() int
	Height() int
}
