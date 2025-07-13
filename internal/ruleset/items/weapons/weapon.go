package weapon

import (
	"fmt"

	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/expression"
	"anvil/internal/ruleset/actions/basic"
	"anvil/internal/tag"
)

type Weapon struct {
	archetype string
	id        string
	name      string
	damage    expression.Expression
	tags      tag.Container
	reach     int
}

func NewWeapon(archetype, id, name string, damage expression.Expression, weaponTags tag.Container, reach int) *Weapon {
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
	tags := w.tags.Clone()
	return &tags
}

func (w Weapon) OnEquip(a *core.Actor) {
	cost := map[tag.Tag]int{tags.ResourceAction: 1}
	a.AddAction(basic.NewMeleeAction(a, fmt.Sprintf("Attack with %s", w.name), &w, w.reach, w.tags, cost))
}

func (w Weapon) Damage() *expression.Expression {
	dmg := w.damage.Clone()
	return &dmg
}
