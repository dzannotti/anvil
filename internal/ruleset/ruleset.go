package ruleset

import (
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/ruleset/base"
)

func newActor(hub *eventbus.Hub, world *core.World, team core.TeamID, pos grid.Position, name string, hitPoints int, attributes stats.Attributes, proficiencies stats.Proficiencies) *core.Actor {
	c := &core.Actor{
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
	world.AddOccupant(pos, c)
	c.AddAction(base.NewAttackAction(c))
	return c
}

func NewPCActor(hub *eventbus.Hub, world *core.World, team core.TeamID, pos grid.Position, name string, hitPoints int, attributes stats.Attributes, proficiencies stats.Proficiencies) *core.Actor {
	c := &core.Actor{
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
	world.AddOccupant(pos, c)
	c.AddAction(base.NewAttackAction(c))
	return c
}
