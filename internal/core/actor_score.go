package core

import (
	"anvil/internal/grid"
)

func (a Actor) BestScoredAction() *ScoredAction {
	return a.BestScoredActionAt(a.Position, func(Action, int) bool { return true })
}

func (a Actor) BestScoredActionAt(pos grid.Position, filter func(Action, int) bool) *ScoredAction {
	return a.BestScoredActionAtWhere(pos, func(_ Action, depth int) bool { return depth < 5 }, 0)
}

func (a Actor) BestScoredActionAtWhere(pos grid.Position, filter func(Action, int) bool, depth int) *ScoredAction {
	var best *ScoredAction
	for _, action := range a.Actions {
		if !filter(action, depth) {
			continue
		}
		for _, pos := range action.ValidPositions(pos) {
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
