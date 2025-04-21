package core

import (
	"sync"

	"anvil/internal/grid"
)

const useAsync = false

func (a Actor) BestScoredAction() *ScoredAction {
	if useAsync {
		return a.bestScoredActionAtAsync(a.Position)
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

func (a *Actor) bestScoredActionAtAsync(pos grid.Position, filter ...func(Action) bool) *ScoredAction {
	scoredCh := make(chan *ScoredAction, 1024)
	var wg sync.WaitGroup
	for _, action := range a.Actions {
		if len(filter) > 0 && filter[0](action) {
			continue
		}
		validPositions := action.ValidPositions(pos)
		for _, validPos := range validPositions {
			wg.Add(1)
			go func(a Action, p grid.Position) {
				defer wg.Done()
				scored := a.ScoreAt(p)
				if scored != nil && scored.Score >= 0.01 {
					scoredCh <- scored
				}
			}(action, validPos)
		}
	}
	go func() {
		wg.Wait()
		close(scoredCh)
	}()

	var best *ScoredAction
	for scored := range scoredCh {
		if best == nil || scored.Score > best.Score {
			best = scored
		}
	}

	return best
}
