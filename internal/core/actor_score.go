package core

import (
	"anvil/internal/grid"
)

func (a Actor) BestScoredAction() *ScoredAction {
	return a.BestScoredActionAt(a.Position)
}

func (a Actor) BestScoredActionAt(pos grid.Position, filter ...func(Action) bool) *ScoredAction {
	var best *ScoredAction
	// TODO: make this concurrent
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

/*

func EvaluateBestTargetSetParallel(action core.Action, src *core.Actor) ([]grid.Position, float32) {
	targetSets := action.ValidTargetSets(src.Position)
	if len(targetSets) == 0 {
		return nil, 0
	}

	results := make(chan scoredResult, len(targetSets))
	var wg sync.WaitGroup

	for _, targets := range targetSets {
		wg.Add(1)
		go func(targets []grid.Position) {
			defer wg.Done()
			score := action.ScoreAt(targets)
			if score >= 0.01 { // skip junk scores
				results <- scoredResult{targets, score}
			}
		}(targets)
	}

	wg.Wait()
	close(results)

	var bestTargets []grid.Position
	var bestScore float32

	for res := range results {
		if res.score > bestScore {
			bestScore = res.score
			bestTargets = res.targets
		}
	}

	return bestTargets, bestScore
}
*/
