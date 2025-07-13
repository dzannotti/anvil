package metrics

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
)

type PositioningMetric struct{}

func (p PositioningMetric) Evaluate(world *core.World, actor *core.Actor, action core.Action, target grid.Position, affected []grid.Position) map[string]int {
	survivalThreat := 0
	enemyProximity := 0
	movementEfficiency := 0
	
	// Get all hostile actors in the encounter
	if actor.Encounter == nil {
		return map[string]int{
			"survival_threat":      0,
			"enemy_proximity":      0,
			"movement_efficiency":  0,
		}
	}
	
	hostileActors := actor.Encounter.HostileActors(actor)
	
	// Calculate survival threat based on enemies that could threaten our current position
	actorPosition := actor.Position
	
	// Step 9: Check if this is a movement action and calculate movement efficiency
	actionTags := action.Tags()
	if actionTags.HasTag(tags.Move) || actionTags.HasTag(tags.Dash) {
		movementEfficiency = calculateMovementEfficiency(world, actor, target, hostileActors)
		// For movement actions, evaluate threat at the TARGET position, not current position
		actorPosition = target
	}
	
	for _, enemy := range hostileActors {
		if enemy.IsDead() {
			continue
		}
		
		// Calculate distance to enemy
		distance := calculateDistance(actorPosition, enemy.Position)
		
		// Enemies within melee range (1 tile) are high threat
		if distance <= 1 {
			survivalThreat += 30
			enemyProximity += 40
		} else if distance <= 3 {
			// Enemies within short range are medium threat
			survivalThreat += 15
			enemyProximity += 20
		} else if distance <= 6 {
			// Enemies within medium range are low threat
			survivalThreat += 5
			enemyProximity += 10
		}
		
		// Additional threat for enemies that can likely reach us next turn
		// Assume most enemies can move + attack within 6 tiles
		if distance <= 6 {
			// Estimate damage this enemy could do to us
			// This is simplified - real implementation would check enemy's actions
			estimatedEnemyDamage := 8 // Rough estimate based on zombie slam
			
			// Higher threat if we're already low on health
			healthRatio := float32(actor.HitPoints) / float32(actor.MaxHitPoints)
			if healthRatio < 0.5 {
				estimatedEnemyDamage = int(float32(estimatedEnemyDamage) * 1.5)
			}
			
			survivalThreat += estimatedEnemyDamage
		}
	}
	
	return map[string]int{
		"survival_threat":      -survivalThreat, // Negative because threat is bad
		"enemy_proximity":      -enemyProximity,  // Negative because being close to enemies is usually bad
		"movement_efficiency":  movementEfficiency, // Can be positive or negative based on tactical value
	}
}

func calculateMovementEfficiency(world *core.World, actor *core.Actor, targetPos grid.Position, hostileActors []*core.Actor) int {
	// Step 9: Calculate movement efficiency based on tactical improvements
	efficiency := 0
	currentPos := actor.Position
	
	// Step 9: Prevent infinite movement loops - penalize staying in same position
	if targetPos == currentPos {
		efficiency -= 5 // Small penalty for not moving when using a movement action
	}
	
	// 1. Calculate movement cost (negative efficiency for longer moves)
	distance := calculateDistance(currentPos, targetPos)
	efficiency -= distance * 2 // 2 points per tile moved
	
	// 2. Evaluate tactical improvements from movement
	
	// Compare threat levels: current position vs target position
	currentThreat := calculateThreatAtPosition(world, currentPos, hostileActors)
	targetThreat := calculateThreatAtPosition(world, targetPos, hostileActors)
	
	// Bonus for moving to safer position
	if targetThreat < currentThreat {
		safetyImprovement := (currentThreat - targetThreat) / 2
		efficiency += safetyImprovement
	}
	
	// 3. Evaluate engagement potential
	if len(hostileActors) > 0 {
		// Find nearest enemy for engagement calculations
		nearestEnemy := findNearestEnemy(currentPos, hostileActors)
		if nearestEnemy != nil {
			currentDistanceToEnemy := calculateDistance(currentPos, nearestEnemy.Position)
			targetDistanceToEnemy := calculateDistance(targetPos, nearestEnemy.Position)
			
			// Slight bonus for getting into better attack range (but not too close)
			if currentDistanceToEnemy > 3 && targetDistanceToEnemy >= 1 && targetDistanceToEnemy <= 2 {
				efficiency += 15 // Good attack position
			} else if currentDistanceToEnemy <= 1 && targetDistanceToEnemy > 1 && targetDistanceToEnemy <= 3 {
				efficiency += 10 // Good retreat to ranged position
			}
		}
	}
	
	// 4. Penalty for moving to obviously bad positions
	if targetPos.X <= 0 || targetPos.Y <= 0 || 
	   targetPos.X >= world.Width()-1 || targetPos.Y >= world.Height()-1 {
		efficiency -= 20 // Avoid corners/edges when possible
	}
	
	return efficiency
}

func calculateThreatAtPosition(world *core.World, pos grid.Position, hostileActors []*core.Actor) int {
	threat := 0
	
	for _, enemy := range hostileActors {
		if enemy.IsDead() {
			continue
		}
		
		distance := calculateDistance(pos, enemy.Position)
		
		// Immediate threat (enemy can reach us next turn)
		if distance <= 1 {
			threat += 30
		} else if distance <= 3 {
			threat += 15
		} else if distance <= 6 {
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

