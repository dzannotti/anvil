package factories

import (
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	actionsBasic "anvil/internal/ruleset/actions/basic"
	effectsBasic "anvil/internal/ruleset/effects/basic"
)

func newActor(
	dispatcher *eventbus.Dispatcher,
	w *core.World,
	t core.TeamID,
	pos grid.Position,
	name string,
	hitPoints int,
	at stats.Attributes,
	p stats.Proficiencies,
	r core.Resources,
) *core.Actor {
	a := &core.Actor{
		Dispatcher:    dispatcher,
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
	a.AddEffect(effectsBasic.NewAttributeModifierEffect())
	a.AddEffect(effectsBasic.NewProficiencyModifierEffect())
	a.AddEffect(effectsBasic.NewCritEffect())
	a.AddAction(actionsBasic.NewMoveAction(a))
	a.Resources.LongRest()
	a.AddEffect(effectsBasic.NewAttackOfOpportunityEffect())
	return a
}

func NewPCActor(
	dispatcher *eventbus.Dispatcher,
	w *core.World,
	pos grid.Position,
	name string,
	hitPoints int,
	at stats.Attributes,
	p stats.Proficiencies,
	r core.Resources,
) *core.Actor {
	a := newActor(dispatcher, w, core.TeamPlayers, pos, name, hitPoints, at, p, r)
	a.AddEffect(effectsBasic.NewDeathSavingThrowEffect())
	return a
}

func NewNPCActor(
	dispatcher *eventbus.Dispatcher,
	w *core.World,
	pos grid.Position,
	name string,
	hitPoints int,
	at stats.Attributes,
	p stats.Proficiencies,
	r core.Resources,
) *core.Actor {
	a := newActor(dispatcher, w, core.TeamEnemies, pos, name, hitPoints, at, p, r)
	a.AddEffect(effectsBasic.NewDeathEffect())
	a.AddProficiency(tags.NaturalWeapon)
	return a
}
