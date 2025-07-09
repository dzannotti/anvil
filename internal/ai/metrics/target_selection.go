package metrics

import (
	"anvil/internal/core"
	"anvil/internal/grid"
)

type TargetSelectionMetric struct{}

func (t TargetSelectionMetric) Evaluate(world *core.World, actor *core.Actor, action core.Action, target grid.Position, affected []grid.Position) map[string]int {
	threatLevel := 0
	lowHealthBonus := 0
	tacticalValue := 0
	
	// Get the primary target actor
	cell := world.At(target)
	if cell == nil || cell.Occupant() == nil {
		return map[string]int{
			"threat_priority":   0,
			"low_health_bonus": 0,
			"tactical_value":   0,
		}
	}
	
	targetActor := cell.Occupant()
	if targetActor.IsDead() || !actor.IsHostileTo(targetActor) {
		return map[string]int{
			"threat_priority":   0,
			"low_health_bonus": 0,
			"tactical_value":   0,
		}
	}
	
	// 1. Threat Level Assessment
	// Higher threat for enemies that can reach us next turn
	distance := calculateDistance(actor.Position, targetActor.Position)
	
	// Assume most enemies can move + attack within 6 tiles
	if distance <= 6 {
		threatLevel += 30
		
		// Closer enemies are higher priority
		if distance <= 3 {
			threatLevel += 20
		}
		if distance <= 1 {
			threatLevel += 30 // Immediate threat
		}
	}
	
	// Higher threat for enemies with more remaining HP (they'll last longer)
	healthRatio := float32(targetActor.HitPoints) / float32(targetActor.MaxHitPoints)
	if healthRatio > 0.75 {
		threatLevel += 25 // Fresh enemies are dangerous
	} else if healthRatio > 0.5 {
		threatLevel += 15
	}
	
	// 2. Low Health Bonus (prioritize finishing off wounded enemies)
	if healthRatio <= 0.25 {
		lowHealthBonus += 40 // Almost dead - finish them off
	} else if healthRatio <= 0.5 {
		lowHealthBonus += 20 // Wounded - good target
	}
	
	// 3. Tactical Value Assessment
	// Estimate enemy's damage potential (simplified)
	enemyDamageEstimate := estimateEnemyDamage(targetActor)
	tacticalValue += enemyDamageEstimate / 2 // Higher damage enemies are higher priority
	
	// Bonus for enemies that are isolated (easier to focus fire)
	nearbyAllies := countNearbyAllies(world, targetActor, 3)
	if nearbyAllies == 0 {
		tacticalValue += 20 // Isolated target
	} else if nearbyAllies == 1 {
		tacticalValue += 10
	}
	
	// Penalty for heavily armored targets (less efficient)
	targetAC := targetActor.ArmorClass()
	if targetAC.Value >= 16 {
		tacticalValue -= 15 // Hard to hit
	} else if targetAC.Value <= 12 {
		tacticalValue += 10 // Easy to hit
	}
	
	return map[string]int{
		"threat_priority":   threatLevel,
		"low_health_bonus": lowHealthBonus,
		"tactical_value":   tacticalValue,
	}
}

func estimateEnemyDamage(enemy *core.Actor) int {
	// Simplified damage estimation based on actions
	maxDamage := 0
	for _, action := range enemy.Actions {
		avgDamage := action.AverageDamage()
		if avgDamage > maxDamage {
			maxDamage = avgDamage
		}
	}
	
	// If no damage info available, estimate based on level/stats
	if maxDamage == 0 {
		// Rough estimate: assume 1d6 + stat mod per "level"
		estimatedLevel := enemy.MaxHitPoints / 8 // Very rough
		if estimatedLevel < 1 {
			estimatedLevel = 1
		}
		maxDamage = 4 + estimatedLevel // 1d6 average + level
	}
	
	return maxDamage
}

func countNearbyAllies(world *core.World, actor *core.Actor, radius int) int {
	count := 0
	pos := actor.Position
	
	// Check all positions within radius
	for x := pos.X - radius; x <= pos.X + radius; x++ {
		for y := pos.Y - radius; y <= pos.Y + radius; y++ {
			checkPos := grid.Position{X: x, Y: y}
			
			// Skip the actor's own position
			if checkPos == pos {
				continue
			}
			
			// Check if position is within actual radius (not just square)
			if calculateDistance(pos, checkPos) > radius {
				continue
			}
			
			cell := world.At(checkPos)
			if cell == nil || cell.Occupant() == nil {
				continue
			}
			
			occupant := cell.Occupant()
			if occupant.IsDead() {
				continue
			}
			
			// Count allies (same team or not hostile)
			if !actor.IsHostileTo(occupant) {
				count++
			}
		}
	}
	
	return count
}