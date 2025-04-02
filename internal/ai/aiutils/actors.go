package aiutils

import (
	"anvil/internal/core"
)

func BestAIAction(src *core.Actor, action core.Action) *core.AIAction {
	var best *core.AIAction
	for _, pos := range action.ValidPositions(src.Position) {
		local := action.AIAction(pos)
		if local.Score < 0.01 {
			continue
		}
		if best == nil {
			best = local
			continue
		}
		if local.Score > best.Score {
			best = local
		}
	}
	return best
}

func BestAIChoice(src *core.Actor) *core.AIAction {
	var best *core.AIAction
	for _, a := range src.Actions {
		local := BestAIAction(src, a)
		if local == nil {
			continue
		}
		if local.Score < 0.01 {
			continue
		}
		if best == nil {
			best = local
			continue
		}
		if local.Score > best.Score {
			best = local
		}
	}
	return best
}
