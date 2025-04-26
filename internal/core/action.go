package core

import (
	"anvil/internal/grid"
	"anvil/internal/tag"
)

type Action interface {
	Name() string
	Tags() *tag.Container
	Perform(pos []grid.Position)

	ValidPositions(from grid.Position) []grid.Position
	AffectedPositions(target []grid.Position) []grid.Position
	ScoreAt(pos grid.Position) float32
	TargetCountAt(pos grid.Position) int
}

type ScoredAction struct {
	Position []grid.Position
	Action   Action
	Score    float32
}
