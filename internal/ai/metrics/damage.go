package metrics

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
)

type DamageMetric struct{}

func (d DamageMetric) Evaluate(world *core.World, actor *core.Actor, action core.Action, target grid.Position, affected []grid.Position) map[string]int {
	enemyDamage := 0
	friendlyFire := 0
	killPotential := 0
	
	// Calculate base damage per target
	baseDamage := action.AverageDamage()
	
	// Analyze each affected position
	for _, pos := range affected {
		cell := world.At(pos)
		if cell == nil || cell.Occupant() == nil {
			continue
		}
		
		occupant := cell.Occupant()
		if occupant.IsDead() {
			continue
		}
		
		// Calculate expected damage considering hit chance
		expectedDamage := calculateExpectedDamage(actor, action, occupant, baseDamage)
		
		if actor.IsHostileTo(occupant) {
			// Enemy damage
			enemyDamage += expectedDamage
			
			// Kill potential: bonus points if this could kill the target
			if expectedDamage >= occupant.HitPoints {
				killPotential += 50 // Bonus for potential kills
			}
		} else {
			// Friendly fire damage (negative)
			friendlyFire += expectedDamage
		}
	}
	
	return map[string]int{
		"damage_enemy":   enemyDamage,
		"friendly_fire":  -friendlyFire, // Negative because it's bad
		"kill_potential": killPotential,
	}
}

func calculateExpectedDamage(attacker *core.Actor, action core.Action, target *core.Actor, baseDamage int) int {
	// For spell attacks and saving throws, assume 75% success rate for simplicity
	if action.Tags().HasTag(tags.Spell) {
		return int(float32(baseDamage) * 0.75)
	}
	
	// For weapon attacks, calculate hit probability based on attack bonus vs AC
	if action.Tags().HasTag(tags.Attack) {
		// Simplified hit calculation: assume +5 attack bonus vs target AC
		// Real implementation would extract actual attack bonus from action/actor
		attackBonus := 5 // Simplified - could be extracted from actor stats
		targetAC := target.ArmorClass()
		
		// Hit on d20 + attackBonus >= targetAC
		// So we need to roll (targetAC - attackBonus) or higher on d20
		neededRoll := targetAC.Value - attackBonus
		
		var hitChance float32
		if neededRoll <= 1 {
			hitChance = 0.95 // Almost always hit (only fail on nat 1)
		} else if neededRoll >= 20 {
			hitChance = 0.05 // Only hit on nat 20
		} else {
			hitChance = float32(21-neededRoll) / 20.0
		}
		
		return int(float32(baseDamage) * hitChance)
	}
	
	// For other action types, assume 100% success
	return baseDamage
}