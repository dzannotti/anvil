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
) int {
	damage := action.AverageDamage()
	if damage == 0 {
		return 0
	}
	targets := targetsAffected(world, affected)
	los := make([]*core.Actor, 0, len(targets))
	for _, t := range targets {
		if world.HasLineOfSight(actor.Position, t.Position) {
			los = append(los, t)
		}
	}
	if len(los) == 0 {
		return 0
	}
	score := BaseDamageScore
	for _, t := range los {
		score += mathi.Min(damage, t.HitPoints)
	}
	return score
}
