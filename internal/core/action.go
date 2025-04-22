package core

import (
	"anvil/internal/grid"
	"anvil/internal/tag"
)

type Action interface {
	Name() string
	ScoreAt(pos grid.Position) *ScoredAction
	Perform(pos []grid.Position)
	ValidPositions(from grid.Position) []grid.Position
	Tags() *tag.Container
	TargetCountAt(pos grid.Position) int
	AffectedPositions(target []grid.Position) []grid.Position
}

type ScoredAction struct {
	Position []grid.Position
	Action   Action
	Score    float32
}
