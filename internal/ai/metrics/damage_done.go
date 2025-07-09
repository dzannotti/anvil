package metrics

import (
	"anvil/internal/core"
	"anvil/internal/grid"
	"anvil/internal/mathi"
)

const BaseDamageScore = 20

type DamageDone struct{}

func (d DamageDone) Evaluate(
	world *core.World,
	actor *core.Actor,
	action core.Action,
	_ grid.Position,
	affected []grid.Position,
) map[string]int {
	damage := action.AverageDamage()
	if damage == 0 {
		return map[string]int{}
	}
	targets := targetsAffected(world, affected)
	los := make([]*core.Actor, 0, len(targets))
	for _, t := range targets {
		if world.HasLineOfSight(actor.Position, t.Position) {
			los = append(los, t)
		}
	}
	if len(los) == 0 {
		return map[string]int{}
	}
	
	totalDamage := 0
	killPotential := 0
	threatElimination := 0
	aoeEfficiency := 0
	overkillWaste := 0
	
	for _, t := range los {
		actualDamage := mathi.Min(damage, t.HitPoints)
		totalDamage += actualDamage
		
		// Kill potential - higher score for enemies we can finish off
		if damage >= t.HitPoints {
			killPotential += 15
		}
		
		// Threat elimination - higher score for dangerous enemies
		if t.HitPoints > 0 {
			threatLevel := mathi.Min(t.MaxHitPoints/4, 10) // Estimate threat by max HP
			threatElimination += threatLevel
		}
		
		// AOE efficiency - bonus for hitting multiple targets
		if len(los) > 1 {
			aoeEfficiency += 5
		}
		
		// Overkill waste - penalty for excessive damage
		if damage > t.HitPoints {
			overkillWaste -= (damage - t.HitPoints) / 2
		}
	}
	
	return map[string]int{
		"damage_dealt":        BaseDamageScore + totalDamage,
		"kill_potential":      killPotential,
		"threat_elimination":  threatElimination,
		"aoe_efficiency":      aoeEfficiency,
		"overkill_waste":      overkillWaste,
	}
}
