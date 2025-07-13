package ai

import (
	"anvil/internal/core"
	"anvil/internal/grid"
	"fmt"
)

type Weights struct {
	Weights map[string]float32
}

type ActionTargetEvaluation struct {
	Action         core.Action
	Target         grid.Position
	Affected       []grid.Position
	RawMetrics     map[string]int
	WeightedScores map[string]int
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
	archetype := getArchetypeName(weights)
	logAIDecision(fmt.Sprintf("\n=== %s (%s) Turn ===", actor.Name, archetype))

	bestEvaluation := findBestAction(state.World, actor, state.Encounter, weights)
	if bestEvaluation != nil && bestEvaluation.FinalScore > 0 {
		logAIDecision(fmt.Sprintf("✅ PRIMARY ACTION: %s -> %v (Score: %d)",
			bestEvaluation.Action.Name(), bestEvaluation.Target, bestEvaluation.FinalScore))
		logActionBreakdown(bestEvaluation)

		executeAction(state, actor, bestEvaluation, originalPos)
		return
	}

	logAIDecision(fmt.Sprintf("⚠️  NO GOOD PRIMARY ACTIONS (best score: %d) - Using fallback",
		getBestScore(bestEvaluation)))

	fallbackAction := findFallbackAction(state.World, actor, state.Encounter, weights)
	if fallbackAction != nil {
		logAIDecision(fmt.Sprintf("🔄 FALLBACK ACTION: %s -> %v (Score: %d)",
			fallbackAction.Action.Name(), fallbackAction.Target, fallbackAction.FinalScore))
		logActionBreakdown(fallbackAction)

		executeAction(state, actor, fallbackAction, originalPos)
		return
	}

	logAIDecision("❌ NO ACTIONS AVAILABLE - Skipping turn")
}

func executeAction(state *core.GameState, actor *core.Actor, evaluation *ActionTargetEvaluation, originalPos grid.Position) {
	if evaluation.Position != actor.Position && len(evaluation.Movement) > 0 {
		state.World.RemoveOccupant(actor.Position, actor)
		actor.Position = evaluation.Position
		state.World.AddOccupant(evaluation.Position, actor)
	}

	evaluation.Action.Perform([]grid.Position{evaluation.Target})

	if actor.Name == "Cedric" && evaluation.Position != originalPos {
		state.World.RemoveOccupant(actor.Position, actor)
		actor.Position = originalPos
		state.World.AddOccupant(originalPos, actor)
	}
}

func getBestScore(evaluation *ActionTargetEvaluation) int {
	if evaluation != nil {
		return evaluation.FinalScore
	}
	return -999
}
