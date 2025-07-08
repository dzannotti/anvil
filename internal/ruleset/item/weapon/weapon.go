package weapon

import (
	"fmt"

	"anvil/internal/core"
	"anvil/internal/expression"
	"anvil/internal/ruleset/base"
	"anvil/internal/tag"
)

type Weapon struct {
	archetype    string
	id           string
	name         string
	damageSource core.DamageSource
	tags         tag.Container
	reach        int
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
	a.AddAction(base.NewAttackAction(a, fmt.Sprintf("Attack with %s", w.name), w.damageSource, w.reach, w.tags))
}

// Implement DamageSource interface
func (w Weapon) Damage() *expression.Expression {
	return w.damageSource.Damage()
}
