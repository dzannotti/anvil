package definition

import "anvil/internal/grid"

type WorldCell interface {
	Position() grid.Position
	Occupant() (Creature, error)
	AddOccupant(Creature)
	RemoveOccupant(Creature)
	IsOccupied() bool
}

type World interface {
	At(grid.Position) (WorldCell, bool)
	CreaturesInRange(grid.Position, int) []Creature
}
