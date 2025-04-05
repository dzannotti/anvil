package core

import (
	"anvil/internal/grid"
	"anvil/internal/tag"
)

type Action interface {
	Name() string
	AIAction(pos grid.Position) *AIAction
	Perform(pos []grid.Position)
	ValidPositions(from grid.Position) []grid.Position
	Tags() tag.Container
}
