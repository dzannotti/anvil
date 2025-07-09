package metrics

import (
	"anvil/internal/core"
	"anvil/internal/grid"
)

type AIMetric interface {
	Evaluate(world *core.World, actor *core.Actor, action core.Action, pos grid.Position, affected []grid.Position) map[string]int
}

var Default = []AIMetric{
	DamageDone{},
	FriendlyFire{},
	Movement{},
	Plan{},
}

func targetsAffected(world *core.World, pos []grid.Position) []*core.Actor {
	targets := make([]*core.Actor, 0, len(pos))
	for _, p := range pos {
		actor := world.ActorAt(p)
		if actor != nil {
			targets = append(targets, actor)
		}
	}
	return targets
}

func enemiesAffected(world *core.World, actor *core.Actor, pos []grid.Position) []*core.Actor {
	targets := targetsAffected(world, pos)
	enemies := make([]*core.Actor, 0, len(targets))
	for _, t := range targets {
		if actor.IsHostileTo(t) {
			enemies = append(enemies, t)
		}
	}
	return enemies
}

func friendliesAffected(world *core.World, actor *core.Actor, pos []grid.Position) []*core.Actor {
	targets := targetsAffected(world, pos)
	friendlies := make([]*core.Actor, 0, len(targets))
	for _, t := range targets {
		if !actor.IsHostileTo(t) {
			friendlies = append(friendlies, t)
		}
	}
	return friendlies
}
