package metrics

import (
	"anvil/internal/core"
	"anvil/internal/grid"
)

type TargetSelectionMetric struct{}

func (t TargetSelectionMetric) Evaluate(world *core.World, actor *core.Actor, _ core.Action, target grid.Position, _ []grid.Position) map[string]int {
	targetActor := getValidTargetActor(world, actor, target)
	if targetActor == nil {
		return emptyTargetMetrics()
	}

	threatLevel := calculateThreatLevel(actor, targetActor)
	lowHealthBonus := calculateLowHealthBonus(targetActor)
	tacticalValue := calculateTacticalValue(world, targetActor)

	return map[string]int{
		"threat_priority":  threatLevel,
		"low_health_bonus": lowHealthBonus,
		"tactical_value":   tacticalValue,
	}
}

func getValidTargetActor(world *core.World, actor *core.Actor, target grid.Position) *core.Actor {
	cell := world.At(target)
	if cell == nil || cell.Occupant() == nil {
		return nil
	}

	targetActor := cell.Occupant()
	if targetActor.IsDead() || !actor.IsHostileTo(targetActor) {
		return nil
	}

	return targetActor
}

func emptyTargetMetrics() map[string]int {
	return map[string]int{
		"threat_priority":  0,
		"low_health_bonus": 0,
		"tactical_value":   0,
	}
}

func calculateThreatLevel(actor *core.Actor, targetActor *core.Actor) int {
	threatLevel := 0
	distance := calculateDistance(actor.Position, targetActor.Position)

	threatLevel += calculateDistanceThreat(distance)
	threatLevel += calculateHealthThreat(targetActor)

	return threatLevel
}

func calculateDistanceThreat(distance int) int {
	if distance > 6 {
		return 0
	}

	threat := 30
	if distance <= 3 {
		threat += 20
	}
	if distance <= 1 {
		threat += 30
	}
	return threat
}

func calculateHealthThreat(targetActor *core.Actor) int {
	healthRatio := float32(targetActor.HitPoints) / float32(targetActor.MaxHitPoints)
	switch {
	case healthRatio > 0.75:
		return 25
	case healthRatio > 0.5:
		return 15
	default:
		return 0
	}
}

func calculateLowHealthBonus(targetActor *core.Actor) int {
	healthRatio := float32(targetActor.HitPoints) / float32(targetActor.MaxHitPoints)
	switch {
	case healthRatio <= 0.25:
		return 40
	case healthRatio <= 0.5:
		return 20
	default:
		return 0
	}
}

func calculateTacticalValue(world *core.World, targetActor *core.Actor) int {
	tacticalValue := estimateEnemyDamage(targetActor) / 2
	tacticalValue += calculateIsolationBonus(world, targetActor)
	tacticalValue += calculateArmorPenalty(targetActor)
	return tacticalValue
}

func calculateIsolationBonus(world *core.World, targetActor *core.Actor) int {
	nearbyAllies := countNearbyAllies(world, targetActor, 3)
	switch nearbyAllies {
	case 0:
		return 20
	case 1:
		return 10
	default:
		return 0
	}
}

func calculateArmorPenalty(targetActor *core.Actor) int {
	targetAC := targetActor.ArmorClass()
	switch {
	case targetAC.Value >= 16:
		return -15
	case targetAC.Value <= 12:
		return 10
	default:
		return 0
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
	for x := pos.X - radius; x <= pos.X+radius; x++ {
		for y := pos.Y - radius; y <= pos.Y+radius; y++ {
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
