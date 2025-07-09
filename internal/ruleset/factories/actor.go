package factories

import (
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
)

// RegistryReader provides read-only access to the registry
type RegistryReader interface {
	NewAction(archetype string, owner *core.Actor, options map[string]interface{}) core.Action
	NewEffect(archetype string, options map[string]interface{}) *core.Effect
	NewItem(archetype string, options map[string]interface{}) core.Item
	NewCreature(archetype string, options map[string]interface{}) *core.Actor
	HasAction(archetype string) bool
	HasEffect(archetype string) bool
	HasItem(archetype string) bool
	HasCreature(archetype string) bool
}

func newActor(
	registry RegistryReader,
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

	// Add basic effects
	a.AddEffect(registry.NewEffect("attribute-modifier", nil))
	a.AddEffect(registry.NewEffect("proficiency-modifier", nil))
	a.AddEffect(registry.NewEffect("critical", nil))
	a.AddEffect(registry.NewEffect("attack-of-opportunity", nil))

	// Add basic actions
	a.AddAction(registry.NewAction("move", a, nil))

	a.Resources.LongRest()
	return a
}

func NewPCActor(
	registry RegistryReader,
	dispatcher *eventbus.Dispatcher,
	w *core.World,
	pos grid.Position,
	name string,
	hitPoints int,
	at stats.Attributes,
	p stats.Proficiencies,
	r core.Resources,
) *core.Actor {
	a := newActor(registry, dispatcher, w, core.TeamPlayers, pos, name, hitPoints, at, p, r)
	a.AddEffect(registry.NewEffect("death-saving-throw", nil))
	return a
}

func NewNPCActor(
	registry RegistryReader,
	dispatcher *eventbus.Dispatcher,
	w *core.World,
	pos grid.Position,
	name string,
	hitPoints int,
	at stats.Attributes,
	p stats.Proficiencies,
	r core.Resources,
) *core.Actor {
	a := newActor(registry, dispatcher, w, core.TeamEnemies, pos, name, hitPoints, at, p, r)
	a.AddEffect(registry.NewEffect("death", nil))
	a.AddProficiency(tags.NaturalWeapon)
	return a
}
