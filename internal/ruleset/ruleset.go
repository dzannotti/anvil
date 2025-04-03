package ruleset

import (
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/ruleset/base"
)

func newActor(hub *eventbus.Hub, world *core.World, team core.TeamID, pos grid.Position, name string, hitPoints int, attributes stats.Attributes, proficiencies stats.Proficiencies) *core.Actor {
	a := &core.Actor{
		Log:           hub,
		Position:      pos,
		World:         world,
		Name:          name,
		Team:          team,
		HitPoints:     hitPoints,
		MaxHitPoints:  hitPoints,
		Attributes:    attributes,
		Proficiencies: proficiencies,
	}
	world.AddOccupant(pos, a)
	a.AddAction(base.NewAttackAction(a))
	a.AddEffect(base.NewDeathEffect(a))
	return a
}

func NewPCActor(hub *eventbus.Hub, world *core.World, pos grid.Position, name string, hitPoints int, attributes stats.Attributes, proficiencies stats.Proficiencies) *core.Actor {
	c := newActor(hub, world, core.TeamPlayers, pos, name, hitPoints, attributes, proficiencies)
	return c
}

func NewNPCActor(hub *eventbus.Hub, world *core.World, pos grid.Position, name string, hitPoints int, attributes stats.Attributes, proficiencies stats.Proficiencies) *core.Actor {
	c := newActor(hub, world, core.TeamEnemies, pos, name, hitPoints, attributes, proficiencies)
	return c
}
