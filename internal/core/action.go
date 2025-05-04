package core

import (
	"anvil/internal/grid"
	"anvil/internal/tag"
)

type Action interface {
	Name() string
	Tags() *tag.Container
	Perform(pos []grid.Position, commitCost bool)

	ValidPositions(from grid.Position) []grid.Position
	AffectedPositions(target []grid.Position) []grid.Position
	AverageDamage() int
}

type ScoredAction struct {
	Position []grid.Position
	Action   Action
	Score    float32
}
