package ai

import (
	"anvil/internal/ai/metrics"
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
)

func findOptimalCastingPosition(world *core.World, actor *core.Actor, action core.Action, target grid.Position, weights *Weights) (grid.Position, []grid.Position) {
	actionTags := action.Tags()

	if actionTags.HasTag(tags.Move) || actionTags.HasTag(tags.Dash) {
		return actor.Position, []grid.Position{target}
	}

	candidatePositions := findCastingPositions(world, actor, action, target)
	if len(candidatePositions) == 0 {
		return actor.Position, []grid.Position{}
	}

	bestPosition := actor.Position
	bestScore := -1000
	bestMovement := []grid.Position{}
	originalPosition := actor.Position

	for _, pos := range candidatePositions {
		actor.Position = pos
		safetyScore := evaluatePositionSafety(world, actor, action, target, weights)
		movementCost := calculateDistance(originalPosition, pos) * 5
		totalScore := safetyScore - movementCost

		if totalScore > bestScore {
			bestScore = totalScore
			bestPosition = pos
			if pos != originalPosition {
				bestMovement = []grid.Position{pos}
			} else {
				bestMovement = []grid.Position{}
			}
		}
	}

	actor.Position = originalPosition
	return bestPosition, bestMovement
}

func findCastingPositions(_ *core.World, actor *core.Actor, action core.Action, target grid.Position) []grid.Position {
	var positions []grid.Position
	movementRange := 6
	startX := actor.Position.X - movementRange
	endX := actor.Position.X + movementRange
	startY := actor.Position.Y - movementRange
	endY := actor.Position.Y + movementRange

	for x := startX; x <= endX; x++ {
		for y := startY; y <= endY; y++ {
			pos := grid.Position{X: x, Y: y}

			if !isValidCastingPosition(actor, pos, action, target) {
				continue
			}

			positions = append(positions, pos)
		}
	}

	return positions
}

func isValidCastingPosition(actor *core.Actor, pos grid.Position, action core.Action, target grid.Position) bool {
	if !isPositionInBounds(actor.World, pos) {
		return false
	}

	if !isPositionAccessible(actor.World, actor, pos) {
		return false
	}

	return canActionHitTargetFromPosition(actor, action, pos, target)
}

func isPositionInBounds(world *core.World, pos grid.Position) bool {
	return pos.X >= 0 && pos.Y >= 0 && pos.X < world.Width() && pos.Y < world.Height()
}

func isPositionAccessible(world *core.World, actor *core.Actor, pos grid.Position) bool {
	cell := world.At(pos)
	if cell == nil {
		return false
	}

	if cell.Occupant() != nil && cell.Occupant() != actor {
		return false
	}

	return cell.Tile != core.Wall
}

func canActionHitTargetFromPosition(actor *core.Actor, action core.Action, pos, target grid.Position) bool {
	originalPos := actor.Position
	actor.Position = pos

	validPositions := action.ValidPositions(pos)
	canHit := positionContainsTarget(validPositions, target)

	actor.Position = originalPos
	return canHit
}

func positionContainsTarget(validPositions []grid.Position, target grid.Position) bool {
	for _, validPos := range validPositions {
		if validPos == target {
			return true
		}
	}
	return false
}

func evaluatePositionSafety(world *core.World, actor *core.Actor, action core.Action, target grid.Position, weights *Weights) int {
	affected := action.AffectedPositions([]grid.Position{target})
	positioningMetric := metrics.PositioningMetric{}
	safetyResults := positioningMetric.Evaluate(world, actor, action, target, affected)

	// Convert positioning metrics to structured scores and apply weights
	rawScores := mapToScores(safetyResults)
	weightedScores := rawScores.ApplyWeights(weights)
	
	// For position safety, we only care about positioning-related scores
	return weightedScores.SurvivalThreat + weightedScores.EnemyProximity + weightedScores.MovementEfficiency
}
