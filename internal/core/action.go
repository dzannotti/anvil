package core

import (
	"anvil/internal/grid"
	"anvil/internal/tag"
)

type Action interface {
	Name() string
	Tags() *tag.Container
	Perform(pos []grid.Position)

	ScoreAt(pos grid.Position) float32
	ValidPositions(from grid.Position) []grid.Position
	TargetCountAt(pos grid.Position) int
	AffectedPositions(target []grid.Position) []grid.Position
}

type ScoredAction struct {
	Position []grid.Position
	Action   Action
	Score    float32
}
