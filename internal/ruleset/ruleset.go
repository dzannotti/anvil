package ruleset

import (
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/ruleset/base"
)

func newActor(h *eventbus.Hub, w *core.World, t core.TeamID, pos grid.Position, name string, hitPoints int, at stats.Attributes, p stats.Proficiencies, r core.Resources) *core.Actor {
	a := &core.Actor{
		Log:           h,
		Position:      pos,
		World:         w,
		Name:          name,
		Team:          t,
		HitPoints:     hitPoints,
		MaxHitPoints:  hitPoints,
		Attributes:    at,
		Proficiencies: p,
		Resources:     r,
	}
	w.AddOccupant(pos, a)
	a.AddEffect(base.NewDeathEffect())
	a.AddEffect(base.NewAttributeModifierEffect())
	a.AddEffect(base.NewProficiencyModifierEffect())
	a.AddEffect(base.NewCritEffect())
	a.Resources.Reset()
	return a
}

func NewPCActor(h *eventbus.Hub, w *core.World, pos grid.Position, name string, hitPoints int, at stats.Attributes, p stats.Proficiencies, r core.Resources) *core.Actor {
	a := newActor(h, w, core.TeamPlayers, pos, name, hitPoints, at, p, r)
	return a
}

func NewNPCActor(h *eventbus.Hub, w *core.World, pos grid.Position, name string, hitPoints int, at stats.Attributes, p stats.Proficiencies, r core.Resources) *core.Actor {
	a := newActor(h, w, core.TeamEnemies, pos, name, hitPoints, at, p, r)
	a.AddProficiency(tags.NaturalWeapon)
	return a
}
