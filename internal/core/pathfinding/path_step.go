package pathfinding

import "anvil/internal/grid"

type PathStep struct {
	Position grid.Position
	GCost    float64 // Cumulative cost from start
	StepCost float64 // Cost for this step
	Distance int     // Distance from start (in steps)
}