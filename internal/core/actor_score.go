package core

import (
	"anvil/internal/grid"
)

func (a Actor) BestScoredAction() *ScoredAction {
	if !a.CanAct() {
		return nil
	}
	return a.bestScoredActionAt(a.Position)
}

func (a Actor) bestScoredActionAt(pos grid.Position, filter ...func(Action) bool) *ScoredAction {
	var best *ScoredAction
	for _, action := range a.Actions {
		if len(filter) > 0 && filter[0](action) {
			continue
		}
		for _, pos := range action.ValidPositions(pos) {
			score := action.ScoreAt(pos)
			if score < 0.01 {
				continue
			}
			if best == nil || score > best.Score {
				best = &ScoredAction{
					Score:    score,
					Action:   action,
					Position: []grid.Position{pos},
				}
			}
		}
	}
	return best
}
