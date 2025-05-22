package metrics

import (
	"anvil/internal/core"
	"anvil/internal/grid"
)

type FriendlyFire struct{}

const FriendlyFireMultiplier = 3

func (d FriendlyFire) Evaluate(
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
	friendlies := friendliesAffected(world, actor, affected)
	if len(friendlies) == 0 {
		return 0
	}
	return -damage*len(friendlies)*FriendlyFireMultiplier - BaseDamageScore
}
