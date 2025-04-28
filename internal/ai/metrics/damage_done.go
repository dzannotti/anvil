package metrics

import (
	"anvil/internal/core"
	"anvil/internal/grid"
)

const BaseDamageScore = 20

type DamageDone struct{}

func (d DamageDone) Evaluate(world *core.World, actor *core.Actor, action core.Action, pos grid.Position, affected []grid.Position) int {
	damage := action.AverageDamage()
	if damage == 0 {
		return 0
	}
	targets := targetsAffected(world, affected)
	return BaseDamageScore + damage*len(targets)
}
