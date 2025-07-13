package metrics

import (
	"anvil/internal/core"
	"anvil/internal/grid"
)

type AIMetric interface {
	Evaluate(world *core.World, actor *core.Actor, action core.Action, pos grid.Position, affected []grid.Position) map[string]int
}

// Default contains the standard AI metrics used for decision making.
var Default = []AIMetric{
	DamageMetric{},
	PositioningMetric{},
	TargetSelectionMetric{},
}

func calculateDistance(pos1, pos2 grid.Position) int {
	dx := pos1.X - pos2.X
	dy := pos1.Y - pos2.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}

	// Manhattan distance
	return dx + dy
}
