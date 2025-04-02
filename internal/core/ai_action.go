package core

import "anvil/internal/grid"

type AIAction struct {
	Position []grid.Position
	Action   Action
	Score    float32
}
