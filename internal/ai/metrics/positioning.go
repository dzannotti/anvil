package metrics

import (
	"anvil/internal/core"
	"anvil/internal/grid"
)

type PositioningMetric struct{}

func (p PositioningMetric) Evaluate(world *core.World, actor *core.Actor, action core.Action, target grid.Position, affected []grid.Position) map[string]int {
	survivalThreat := 0
	enemyProximity := 0
	
	// Get all hostile actors in the encounter
	if actor.Encounter == nil {
		return map[string]int{
			"survival_threat": 0,
			"enemy_proximity": 0,
		}
	}
	
	hostileActors := actor.Encounter.HostileActors(actor)
	
	// Calculate survival threat based on enemies that could threaten our current position
	actorPosition := actor.Position
	
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
		"survival_threat": -survivalThreat, // Negative because threat is bad
		"enemy_proximity": -enemyProximity,  // Negative because being close to enemies is usually bad
	}
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