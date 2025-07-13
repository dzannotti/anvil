package ai

import (
	"anvil/internal/core"
	"anvil/internal/grid"
)


type ActionTargetEvaluation struct {
	Action         core.Action
	Target         grid.Position
	Affected       []grid.Position
	RawScores      *Scores
	WeightedScores *WeightedScores
	FinalScore     int
	Position       grid.Position
	Movement       []grid.Position
}

func Play(state *core.GameState, weights *Weights) {
	defer state.Encounter.EndTurn()

	actor := state.Encounter.ActiveActor()
	if !actor.CanAct() || actor.IsDead() {
		return
	}

	originalPos := actor.Position

	bestEvaluation := findBestAction(state.World, actor, state.Encounter, weights)
	if bestEvaluation != nil && bestEvaluation.FinalScore > 0 {
		executeAction(state, actor, bestEvaluation, originalPos)
		return
	}


	fallbackAction := findFallbackAction(state.World, actor, state.Encounter, weights)
	if fallbackAction != nil {
		executeAction(state, actor, fallbackAction, originalPos)
		return
	}

}

func executeAction(state *core.GameState, actor *core.Actor, evaluation *ActionTargetEvaluation, originalPos grid.Position) {
	evaluation.Action.Perform([]grid.Position{evaluation.Target})
}

func getBestScore(evaluation *ActionTargetEvaluation) int {
	if evaluation != nil {
		return evaluation.FinalScore
	}
	return -999
}
