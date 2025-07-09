package weapon

import (
	"fmt"

	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/expression"
	"anvil/internal/ruleset/actions/basic"
	"anvil/internal/tag"
)

type DamageEntry struct {
	Times int
	Sides int
	Kind  tag.Tag
}

type Weapon struct {
	archetype string
	id        string
	name      string
	damage    []DamageEntry
	tags      tag.Container
	reach     int
}

func NewWeapon(archetype, id, name string, damage []DamageEntry, weaponTags tag.Container, reach int) *Weapon {
	return &Weapon{
		archetype: archetype,
		id:        id,
		name:      name,
		damage:    damage,
		tags:      weaponTags,
		reach:     reach,
	}
}

func (w Weapon) Archetype() string {
	return w.archetype
}

func (w Weapon) ID() string {
	return w.id
}

func (w Weapon) Name() string {
	return w.name
}

func (w Weapon) Tags() *tag.Container {
	return &w.tags
}

func (w Weapon) OnEquip(a *core.Actor) {
	cost := map[tag.Tag]int{tags.Action: 1}
	a.AddAction(basic.NewMeleeAction(a, fmt.Sprintf("Attack with %s", w.name), &w, w.reach, w.tags, cost))
}

func (w Weapon) Damage() *expression.Expression {
	if len(w.damage) == 0 {
		// Return empty expression if no damage entries
		return &expression.Expression{}
	}

	// Start with the first damage entry
	first := w.damage[0]
	expr := expression.FromDamageDice(first.Times, first.Sides, w.name, tag.NewContainer(first.Kind))

	// Add remaining damage entries
	for i := 1; i < len(w.damage); i++ {
		entry := w.damage[i]
		expr.AddDamageDice(entry.Times, entry.Sides, w.name, tag.NewContainer(entry.Kind))
	}

	return &expr
}
