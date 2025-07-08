package core

import (
	"anvil/internal/grid"
	"anvil/internal/tag"
)

type Action interface {
	Name() string
	Archetype() string
	ID() string
	Tags() *tag.Container
	Perform(pos []grid.Position)

	ValidPositions(from grid.Position) []grid.Position
	AffectedPositions(target []grid.Position) []grid.Position
	AverageDamage() int
}
