package metrics

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
)

type PositioningMetric struct{}

func (p PositioningMetric) Evaluate(world *core.World, actor *core.Actor, action core.Action, target grid.Position, _ []grid.Position) map[string]int {
	if actor.Encounter == nil {
		return emptyPositioningMetrics()
	}

	hostileActors := actor.Encounter.HostileActors(actor)
	actorPosition, movementEfficiency := determineEvaluationPosition(world, actor, action, target, hostileActors)

	survivalThreat, enemyProximity := calculatePositionThreats(actorPosition, actor, hostileActors)

	return map[string]int{
		"survival_threat":     -survivalThreat,
		"enemy_proximity":     -enemyProximity,
		"movement_efficiency": movementEfficiency,
	}
}

func emptyPositioningMetrics() map[string]int {
	return map[string]int{
		"survival_threat":     0,
		"enemy_proximity":     0,
		"movement_efficiency": 0,
	}
}

func determineEvaluationPosition(world *core.World, actor *core.Actor, action core.Action, target grid.Position, hostileActors []*core.Actor) (grid.Position, int) {
	actionTags := action.Tags()
	if actionTags.HasTag(tags.Move) || actionTags.HasTag(tags.Dash) {
		movementEfficiency := calculateMovementEfficiency(world, actor, target, hostileActors)
		return target, movementEfficiency
	}
	return actor.Position, 0
}

func calculatePositionThreats(actorPosition grid.Position, actor *core.Actor, hostileActors []*core.Actor) (int, int) {
	survivalThreat := 0
	enemyProximity := 0

	for _, enemy := range hostileActors {
		if enemy.IsDead() {
			continue
		}

		distance := calculateDistance(actorPosition, enemy.Position)
		proximityThreat, proximityPenalty := calculateProximityScores(distance)

		survivalThreat += proximityThreat
		enemyProximity += proximityPenalty

		if distance <= 6 {
			survivalThreat += calculateDamageBasedThreat(actor)
		}
	}

	return survivalThreat, enemyProximity
}

func calculateProximityScores(distance int) (int, int) {
	switch {
	case distance <= 1:
		return 30, 40
	case distance <= 3:
		return 15, 20
	case distance <= 6:
		return 5, 10
	default:
		return 0, 0
	}
}

func calculateDamageBasedThreat(actor *core.Actor) int {
	estimatedEnemyDamage := 8
	healthRatio := float32(actor.HitPoints) / float32(actor.MaxHitPoints)
	if healthRatio < 0.5 {
		estimatedEnemyDamage = int(float32(estimatedEnemyDamage) * 1.5)
	}
	return estimatedEnemyDamage
}

func calculateMovementEfficiency(world *core.World, actor *core.Actor, targetPos grid.Position, hostileActors []*core.Actor) int {
	efficiency := calculateBasicMovementCost(actor.Position, targetPos)
	efficiency += calculateSafetyImprovement(world, actor.Position, targetPos, hostileActors)
	efficiency += calculateEngagementBonus(actor.Position, targetPos, hostileActors)
	efficiency += calculatePositionPenalty(world, targetPos)
	return efficiency
}

func calculateBasicMovementCost(currentPos, targetPos grid.Position) int {
	efficiency := 0

	if targetPos == currentPos {
		efficiency -= 5
	}

	distance := calculateDistance(currentPos, targetPos)
	efficiency -= distance * 2

	return efficiency
}

func calculateSafetyImprovement(_ *core.World, currentPos, targetPos grid.Position, hostileActors []*core.Actor) int {
	currentThreat := calculateThreatAtPosition(nil, currentPos, hostileActors)
	targetThreat := calculateThreatAtPosition(nil, targetPos, hostileActors)

	if targetThreat < currentThreat {
		return (currentThreat - targetThreat) / 2
	}
	return 0
}

func calculateEngagementBonus(currentPos, targetPos grid.Position, hostileActors []*core.Actor) int {
	if len(hostileActors) == 0 {
		return 0
	}

	nearestEnemy := findNearestEnemy(currentPos, hostileActors)
	if nearestEnemy == nil {
		return 0
	}

	currentDistanceToEnemy := calculateDistance(currentPos, nearestEnemy.Position)
	targetDistanceToEnemy := calculateDistance(targetPos, nearestEnemy.Position)

	switch {
	case currentDistanceToEnemy > 3 && targetDistanceToEnemy >= 1 && targetDistanceToEnemy <= 2:
		return 15
	case currentDistanceToEnemy <= 1 && targetDistanceToEnemy > 1 && targetDistanceToEnemy <= 3:
		return 10
	default:
		return 0
	}
}

func calculatePositionPenalty(world *core.World, targetPos grid.Position) int {
	if targetPos.X <= 0 || targetPos.Y <= 0 ||
		targetPos.X >= world.Width()-1 || targetPos.Y >= world.Height()-1 {
		return -20
	}
	return 0
}

func calculateThreatAtPosition(_ *core.World, pos grid.Position, hostileActors []*core.Actor) int {
	threat := 0

	for _, enemy := range hostileActors {
		if enemy.IsDead() {
			continue
		}

		distance := calculateDistance(pos, enemy.Position)

		switch {
		case distance <= 1:
			threat += 30
		case distance <= 3:
			threat += 15
		case distance <= 6:
			threat += 5
		}
	}

	return threat
}

func findNearestEnemy(pos grid.Position, hostileActors []*core.Actor) *core.Actor {
	var nearest *core.Actor
	minDistance := 999

	for _, enemy := range hostileActors {
		if enemy.IsDead() {
			continue
		}

		distance := calculateDistance(pos, enemy.Position)
		if distance < minDistance {
			minDistance = distance
			nearest = enemy
		}
	}

	return nearest
}
