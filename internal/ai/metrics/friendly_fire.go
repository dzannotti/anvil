package metrics

import (
	"anvil/internal/core"
	"anvil/internal/grid"
)

type FriendlyFire struct{}

const FriendlyFireMultiplier = 2

func (d FriendlyFire) Evaluate(world *core.World, actor *core.Actor, action core.Action, pos grid.Position, affected []grid.Position) int {
	damage := action.AverageDamage()
	if damage == 0 {
		return 0
	}
	friendlies := friedliesAffected(world, actor, affected)
	return -damage * len(friendlies) * FriendlyFireMultiplier
}
