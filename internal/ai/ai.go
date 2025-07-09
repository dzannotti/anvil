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
			"damage_enemy":      2.0,
			"friendly_fire":     0.5,
			"survival_threat":   0.3,
			"kill_potential":    1.5,
			"enemy_proximity":   0.2,
			"threat_priority":   0.8,  // Berserkers care less about target selection
			"low_health_bonus":  1.8,  // Focus fire on wounded enemies
			"tactical_value":    0.6,
		},
	}
}

func NewDefensiveWeights() *AIWeights {
	return &AIWeights{
		Weights: map[string]float32{
			"damage_enemy":      1.0,
			"friendly_fire":     2.0,
			"survival_threat":   2.0,
			"kill_potential":    0.8,
			"enemy_proximity":   1.8,
			"threat_priority":   1.5,  // Defensive AI prioritizes threats
			"low_health_bonus":  1.0,  // Less focus fire, more threat management
			"tactical_value":    1.3,
		},
	}
}

func NewDefaultWeights() *AIWeights {
	return &AIWeights{
		Weights: map[string]float32{
			"damage_enemy":      1.0,
			"friendly_fire":     1.5,
			"survival_threat":   1.0,
			"kill_potential":    1.2,
			"enemy_proximity":   1.0,
			"threat_priority":   1.2,  // Balanced threat assessment
			"low_health_bonus":  1.4,  // Good focus fire
			"tactical_value":    1.0,
		},
	}
}

func Play(state *core.GameState, weights *AIWeights) {
	defer state.Encounter.EndTurn()
	
	actor := state.Encounter.ActiveActor()
	if !actor.CanAct() || actor.IsDead() {
		return
	}
	
	// Store original position for debugging
	originalPos := actor.Position
	
	// Use new target-centric evaluation flow
	bestEvaluation := findBestAction(state.World, actor, state.Encounter, weights)
	if bestEvaluation != nil {
		// Move to optimal position if needed
		if bestEvaluation.Position != actor.Position && len(bestEvaluation.Movement) > 0 {
			// For now, just teleport to the position. Later we'll implement proper movement
			// Remove from old position and add to new position
			state.World.RemoveOccupant(actor.Position, actor)
			actor.Position = bestEvaluation.Position
			state.World.AddOccupant(bestEvaluation.Position, actor)
		}
		
		bestEvaluation.Action.Perform([]grid.Position{bestEvaluation.Target})
		
		// DEBUG: Force player back to original position to see position optimization in action
		if actor.Name == "Cedric" && bestEvaluation.Position != originalPos {
			state.World.RemoveOccupant(actor.Position, actor)
			actor.Position = originalPos
			state.World.AddOccupant(originalPos, actor)
		}
		
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
		
		// Check if this action can theoretically reach this target from any reachable position
		if canReachTargetFromAnyPosition(world, actor, action, target) {
			potentialTargets = append(potentialTargets, target)
		}
	}
	
	return potentialTargets
}

func simulateActionTarget(world *core.World, actor *core.Actor, action core.Action, target grid.Position, weights *AIWeights) *ActionTargetEvaluation {
	// Find the best position to cast this action from
	bestPosition, movement := findOptimalCastingPosition(world, actor, action, target, weights)
	
	affected := action.AffectedPositions([]grid.Position{target})
	
	// Temporarily move actor to best position for evaluation
	originalPosition := actor.Position
	actor.Position = bestPosition
	
	// Run all metrics to simulate this action-target combination from the optimal position
	rawMetrics := make(map[string]int)
	for _, metric := range metrics.Default {
		metricResults := metric.Evaluate(world, actor, action, target, affected)
		for key, value := range metricResults {
			rawMetrics[key] = value
		}
	}
	
	// Restore original position
	actor.Position = originalPosition
	
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
		Position:       bestPosition,
		Movement:       movement,
	}
}

func canReachTargetFromAnyPosition(world *core.World, actor *core.Actor, action core.Action, target grid.Position) bool {
	// Check if we can reach this target from current position
	validPositions := action.ValidPositions(actor.Position)
	for _, validPos := range validPositions {
		if validPos == target {
			return true
		}
	}
	
	// For now, assume we can reach any target by moving.
	// Later we'll implement proper pathfinding and movement range checking
	return true
}

func findOptimalCastingPosition(world *core.World, actor *core.Actor, action core.Action, target grid.Position, weights *AIWeights) (grid.Position, []grid.Position) {
	// Find all positions from which we can cast this action at the target
	candidatePositions := findCastingPositions(world, actor, action, target)
	
	if len(candidatePositions) == 0 {
		// No valid positions, stay where we are
		return actor.Position, []grid.Position{}
	}
	
	// Evaluate each position for safety and tactical value
	bestPosition := actor.Position
	bestScore := -1000
	bestMovement := []grid.Position{}
	
	originalPosition := actor.Position
	
	for _, pos := range candidatePositions {
		// Temporarily move actor to evaluate this position
		actor.Position = pos
		
		// Calculate safety score for this position
		safetyScore := evaluatePositionSafety(world, actor, action, target, weights)
		
		// Prefer positions closer to current position (less movement cost)
		movementCost := calculateDistance(originalPosition, pos) * 5 // 5 points per tile moved
		totalScore := safetyScore - movementCost
		
		if totalScore > bestScore {
			bestScore = totalScore
			bestPosition = pos
			// For now, movement is just a direct line. Later we'll implement proper pathfinding
			if pos != originalPosition {
				bestMovement = []grid.Position{pos}
			} else {
				bestMovement = []grid.Position{}
			}
		}
	}
	
	// Restore original position
	actor.Position = originalPosition
	
	return bestPosition, bestMovement
}

func findCastingPositions(world *core.World, actor *core.Actor, action core.Action, target grid.Position) []grid.Position {
	var positions []grid.Position
	
	// Check all positions within reasonable movement range (for now, assume 6 tile movement)
	movementRange := 6
	startX := actor.Position.X - movementRange
	endX := actor.Position.X + movementRange
	startY := actor.Position.Y - movementRange
	endY := actor.Position.Y + movementRange
	
	for x := startX; x <= endX; x++ {
		for y := startY; y <= endY; y++ {
			pos := grid.Position{X: x, Y: y}
			
			// Skip if position is out of world bounds
			if pos.X < 0 || pos.Y < 0 || pos.X >= world.Width() || pos.Y >= world.Height() {
				continue
			}
			
			// Skip if position is occupied by someone else
			cell := world.At(pos)
			if cell == nil || (cell.Occupant() != nil && cell.Occupant() != actor) {
				continue
			}
			
			// Skip if position is blocked by terrain
			if cell.Tile == core.Wall {
				continue
			}
			
			// Check if we can cast the action at the target from this position
			// Temporarily move actor to test position
			originalPos := actor.Position
			actor.Position = pos
			
			validPositions := action.ValidPositions(pos)
			canHitTarget := false
			for _, validPos := range validPositions {
				if validPos == target {
					canHitTarget = true
					break
				}
			}
			
			// Restore original position
			actor.Position = originalPos
			
			if canHitTarget {
				positions = append(positions, pos)
			}
		}
	}
	
	return positions
}

func evaluatePositionSafety(world *core.World, actor *core.Actor, action core.Action, target grid.Position, weights *AIWeights) int {
	// Use our existing positioning metric to evaluate safety
	affected := action.AffectedPositions([]grid.Position{target})
	
	positioningMetric := metrics.PositioningMetric{}
	safetyResults := positioningMetric.Evaluate(world, actor, action, target, affected)
	
	// Apply weights to get final safety score
	totalScore := 0
	for metricName, rawValue := range safetyResults {
		if multiplier, exists := weights.Weights[metricName]; exists {
			weightedScore := int(float32(rawValue) * multiplier)
			totalScore += weightedScore
		}
	}
	
	return totalScore
}

func calculateDistance(pos1, pos2 grid.Position) int {
	dx := pos1.X - pos2.X
	dy := pos1.Y - pos2.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	
	// Manhattan distance
	return dx + dy
}