package pathfinding

import (
	"anvil/internal/grid"
	"math"
)

type Result struct {
	Steps     []PathStep
	TotalCost float64
	Found     bool
}

func (r *Result) Positions() []grid.Position {
	if !r.Found {
		return nil
	}
	positions := make([]grid.Position, len(r.Steps))
	for i, step := range r.Steps {
		positions[i] = step.Position
	}
	return positions
}

func (r *Result) Speed() int {
	return int(math.Ceil(r.TotalCost))
}
