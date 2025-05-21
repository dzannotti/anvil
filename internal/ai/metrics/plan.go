package metrics

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
	"anvil/internal/mathi"
	"anvil/internal/tag"
)

const planWeight = float32(0.8)

type Plan struct{}

func (p Plan) Evaluate(
	world *core.World,
	actor *core.Actor,
	action core.Action,
	pos grid.Position,
	_ []grid.Position,
) int {
	if !action.Tags().MatchTag(tags.Move) {
		return 0
	}
	best := 0
	damage := DamageDone{}
	friendly := FriendlyFire{}
	oldPos := actor.Position
	actor.Position = pos
	world.RemoveOccupant(oldPos, actor)
	world.AddOccupant(pos, actor)
	for _, suba := range actor.Actions {
		if suba.Tags().MatchAnyTag(tag.ContainerFromTag(tags.Move, tags.Dash)) {
			continue
		}
		valid := suba.ValidPositions(pos)
		for _, p := range valid {
			affected := suba.AffectedPositions([]grid.Position{p})
			dmg := damage.Evaluate(world, actor, suba, p, affected)
			friendly := friendly.Evaluate(world, actor, suba, p, affected)
			score := dmg + friendly
			best = mathi.Max(int(float32(score)*planWeight), best)
		}
	}
	actor.Position = oldPos
	world.RemoveOccupant(pos, actor)
	world.AddOccupant(oldPos, actor)
	return best
}
