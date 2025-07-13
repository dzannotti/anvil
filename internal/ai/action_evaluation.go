package ai

import (
	"anvil/internal/ai/metrics"
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
	"anvil/internal/tag"
)

func findBestAction(world *core.World, actor *core.Actor, encounter *core.Encounter, weights *Weights) *ActionTargetEvaluation {
	var bestEvaluation *ActionTargetEvaluation
	bestScore := -1000

	for _, action := range actor.Actions {
		if !checkActionFeasibility(action, actor) {
			continue
		}

		targets := findPotentialTargets(world, actor, action, encounter)
		for _, target := range targets {
			evaluation := simulateActionTarget(world, actor, action, target, weights)
			if evaluation.FinalScore > bestScore {
				bestScore = evaluation.FinalScore
				bestEvaluation = evaluation
			}
		}
	}

	return bestEvaluation
}

func findFallbackAction(world *core.World, actor *core.Actor, encounter *core.Encounter, weights *Weights) *ActionTargetEvaluation {
	evaluations := evaluateAllFallbackActions(world, actor, encounter, weights)
	return selectBestFallbackAction(evaluations)
}

func evaluateAllFallbackActions(world *core.World, actor *core.Actor, encounter *core.Encounter, weights *Weights) *fallbackEvaluations {
	evals := &fallbackEvaluations{bestScore: -9999}

	for _, action := range actor.Actions {
		if !checkActionFeasibility(action, actor) {
			continue
		}

		targets := findPotentialTargets(world, actor, action, encounter)
		if len(targets) == 0 {
			continue
		}

		evaluateActionTargets(world, actor, action, targets, weights, evals)
	}

	return evals
}

func evaluateActionTargets(world *core.World, actor *core.Actor, action core.Action, targets []grid.Position, weights *Weights, evals *fallbackEvaluations) {
	actionTags := action.Tags()

	for _, target := range targets {
		evaluation := simulateActionTarget(world, actor, action, target, weights)

		updateBestEvaluation(evaluation, *actionTags, evals)
	}
}

func updateBestEvaluation(evaluation *ActionTargetEvaluation, actionTags tag.Container, evals *fallbackEvaluations) {
	switch {
	case actionTags.HasTag(tags.Move) || actionTags.HasTag(tags.Dash):
		if evals.bestMovement == nil || evaluation.FinalScore > evals.bestMovement.FinalScore {
			evals.bestMovement = evaluation
		}
	case actionTags.HasTag(tags.Dodge):
		if evals.bestDefensive == nil || evaluation.FinalScore > evals.bestDefensive.FinalScore {
			evals.bestDefensive = evaluation
		}
	}

	if evaluation.FinalScore > evals.bestScore {
		evals.bestScore = evaluation.FinalScore
		evals.bestOverall = evaluation
	}
}

func selectBestFallbackAction(evals *fallbackEvaluations) *ActionTargetEvaluation {
	if evals.bestMovement != nil && evals.bestMovement.FinalScore > -50 {
		return evals.bestMovement
	}

	if evals.bestDefensive != nil && evals.bestDefensive.FinalScore > evals.bestScore {
		return evals.bestDefensive
	}

	return evals.bestOverall
}

type fallbackEvaluations struct {
	bestMovement  *ActionTargetEvaluation
	bestDefensive *ActionTargetEvaluation
	bestOverall   *ActionTargetEvaluation
	bestScore     int
}

func simulateActionTarget(world *core.World, actor *core.Actor, action core.Action, target grid.Position, weights *Weights) *ActionTargetEvaluation {
	bestPosition, movement := findOptimalCastingPosition(world, actor, action, target, weights)
	affected := action.AffectedPositions([]grid.Position{target})

	originalPosition := actor.Position
	actor.Position = bestPosition

	rawMetricsMap := make(map[string]int)
	for _, metric := range metrics.Default {
		metricResults := metric.Evaluate(world, actor, action, target, affected)
		for key, value := range metricResults {
			rawMetricsMap[key] = value
		}
	}

	actor.Position = originalPosition

	// Convert map-based metrics to struct-based scores
	rawScores := mapToScores(rawMetricsMap)
	weightedScores := rawScores.ApplyWeights(weights)

	return &ActionTargetEvaluation{
		Action:         action,
		Target:         target,
		Affected:       affected,
		RawScores:      rawScores,
		WeightedScores: weightedScores,
		FinalScore:     weightedScores.Total(),
		Position:       bestPosition,
		Movement:       movement,
	}
}
