package core

import "anvil/internal/grid"

type Action interface {
	Name() string
	AIAction(pos grid.Position) *AIAction
	Perform(pos []grid.Position)
	ValidPositions(from grid.Position) []grid.Position
}
