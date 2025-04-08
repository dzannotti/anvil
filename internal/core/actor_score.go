package core

import (
	"anvil/internal/grid"
)

func (a Actor) BestScoredAction() *ScoredAction {
	return a.BestScoredActionAt(a.Position)
}

func (a Actor) BestScoredActionAt(pos grid.Position, filter ...func(Action) bool) *ScoredAction {
	var best *ScoredAction
	// REMINDER: This cannot be made concurrent if an action uses pathfinding
	for _, action := range a.Actions {
		for _, pos := range action.ValidPositions(pos) {
			if len(filter) > 0 && filter[0](action) {
				continue
			}
			scored := action.ScoreAt(pos)
			if scored == nil || scored.Score < 0.01 {
				continue
			}
			if best == nil || scored.Score > best.Score {
				best = scored
			}
		}
	}

	return best
}
