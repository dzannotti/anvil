package core

import (
	"anvil/internal/grid"
)

func (a Actor) BestScoredAction() *ScoredAction {
	return a.BestScoredActionAt(a.Position)
}

func (a Actor) BestScoredActionAt(pos grid.Position) *ScoredAction {
	var best *ScoredAction

	for _, action := range a.Actions {
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
