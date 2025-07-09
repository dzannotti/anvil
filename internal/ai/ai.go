package ai

import (
	"anvil/internal/ai/metrics"
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
)

type AIWeights struct {
	Weights map[string]float32
}

type ActionTargetEvaluation struct {
	Action       core.Action
	Target       grid.Position
	Affected     []grid.Position
	RawMetrics   map[string]int  // Raw metric values from simulation
	WeightedScores map[string]int  // After applying archetype weights
	FinalScore   int             // Total weighted score
	Position     grid.Position   // Optimal casting position (for now, current position)
	Movement     []grid.Position // Path to position (for now, empty)
}

func NewBerserkerWeights() *AIWeights {
	return &AIWeights{
		Weights: map[string]float32{
			"damage_enemy":     2.0,
			"friendly_fire":    0.5,
			"survival_threat":  0.3,
			"kill_potential":   1.5,
			"enemy_proximity":  0.2,
		},
	}
}

func NewDefensiveWeights() *AIWeights {
	return &AIWeights{
		Weights: map[string]float32{
			"damage_enemy":     1.0,
			"friendly_fire":    2.0,
			"survival_threat":  2.0,
			"kill_potential":   0.8,
			"enemy_proximity":  1.8,
		},
	}
}

func NewDefaultWeights() *AIWeights {
	return &AIWeights{
		Weights: map[string]float32{
			"damage_enemy":     1.0,
			"friendly_fire":    1.5,
			"survival_threat":  1.0,
			"kill_potential":   1.2,
			"enemy_proximity":  1.0,
		},
	}
}

func Play(state *core.GameState, weights *AIWeights) {
	defer state.Encounter.EndTurn()
	
	actor := state.Encounter.ActiveActor()
	if !actor.CanAct() || actor.IsDead() {
		return
	}
	
	// Use new target-centric evaluation flow
	bestEvaluation := findBestAction(state.World, actor, state.Encounter, weights)
	if bestEvaluation != nil {
		bestEvaluation.Action.Perform([]grid.Position{bestEvaluation.Target})
		return
	}
	
	// No valid actions, skip turn
}

func findBestAction(world *core.World, actor *core.Actor, encounter *core.Encounter, weights *AIWeights) *ActionTargetEvaluation {
	var bestEvaluation *ActionTargetEvaluation
	bestScore := -1000
	
	// 1. Check feasible actions
	for _, action := range actor.Actions {
		if !checkActionFeasibility(action, actor) {
			continue
		}
		
		// 2. Find potential targets for this action
		targets := findPotentialTargets(world, actor, action, encounter)
		
		// 3. Simulate action-target combinations
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

func checkActionFeasibility(action core.Action, actor *core.Actor) bool {
	// Check if actor can act (not incapacitated, unconscious, dead, etc.)
	if !actor.CanAct() || actor.IsDead() {
		return false
	}
	
	// Check if actor can afford this action (resources, spell slots, etc.)
	if !action.CanAfford() {
		return false
	}
	
	// Accept all action types: Attack, Spell, Move, Dash, Dodge, Help, Teleport
	actionTags := action.Tags()
	return actionTags.HasTag(tags.Attack) ||
		   actionTags.HasTag(tags.Spell) ||
		   actionTags.HasTag(tags.Move) ||
		   actionTags.HasTag(tags.Dash) ||
		   actionTags.HasTag(tags.Dodge) ||
		   actionTags.HasTag(tags.Help) ||
		   actionTags.HasTag(tags.Teleport)
}

func findPotentialTargets(world *core.World, actor *core.Actor, action core.Action, encounter *core.Encounter) []grid.Position {
	// Get all hostile actors in the encounter (like a human would see)
	hostileActors := encounter.HostileActors(actor)
	
	var potentialTargets []grid.Position
	for _, hostileActor := range hostileActors {
		target := hostileActor.Position
		
		// Check if this action can theoretically reach this target
		// For now, just check if it's in the action's valid positions
		// TODO: Later this will be handled by position optimization
		validPositions := action.ValidPositions(actor.Position)
		for _, validPos := range validPositions {
			if validPos == target {
				potentialTargets = append(potentialTargets, target)
				break
			}
		}
	}
	
	return potentialTargets
}

func simulateActionTarget(world *core.World, actor *core.Actor, action core.Action, target grid.Position, weights *AIWeights) *ActionTargetEvaluation {
	affected := action.AffectedPositions([]grid.Position{target})
	
	// Run all metrics to simulate this action-target combination
	rawMetrics := make(map[string]int)
	for _, metric := range metrics.Default {
		metricResults := metric.Evaluate(world, actor, action, target, affected)
		for key, value := range metricResults {
			rawMetrics[key] = value
		}
	}
	
	// Apply archetype weights to raw metrics
	weightedScores := make(map[string]int)
	totalScore := 0
	for metricName, rawValue := range rawMetrics {
		if multiplier, exists := weights.Weights[metricName]; exists {
			weightedScore := int(float32(rawValue) * multiplier)
			weightedScores[metricName] = weightedScore
			totalScore += weightedScore
		}
	}
	
	return &ActionTargetEvaluation{
		Action:         action,
		Target:         target,
		Affected:       affected,
		RawMetrics:     rawMetrics,
		WeightedScores: weightedScores,
		FinalScore:     totalScore,
		Position:       actor.Position, // For now, always current position
		Movement:       []grid.Position{}, // For now, no movement
	}
}