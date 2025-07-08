package weapon

import (
	"fmt"

	"anvil/internal/core"
	"anvil/internal/expression"
	"anvil/internal/ruleset/base"
	"anvil/internal/tag"
)

type Weapon struct {
	archetype   string
	id          string
	name        string
	damageTimes int
	damageSides int
	damageTags  tag.Container
	tags        tag.Container
	reach       int
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
	a.AddAction(base.NewAttackAction(a, fmt.Sprintf("Attack with %s", w.name), &w, w.reach, w.tags))
}

// Implement DamageSource interface
func (w Weapon) Damage() *expression.Expression {
	expr := expression.FromDamageDice(w.damageTimes, w.damageSides, w.name, w.damageTags)
	return &expr
}

func (w Weapon) DamageTags() *tag.Container {
	return &w.damageTags
}
