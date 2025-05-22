package weapon

import (
	"fmt"

	"anvil/internal/core"
	"anvil/internal/ruleset/base"
	"anvil/internal/tag"
)

type Weapon struct {
	name   string
	damage []core.DamageSource
	tags   tag.Container
	reach  int
}

func (w Weapon) Name() string {
	return w.name
}

func (w Weapon) Tags() *tag.Container {
	return &w.tags
}

func (w Weapon) OnEquip(a *core.Actor) {
	a.AddAction(base.NewAttackAction(a, fmt.Sprintf("Attack with %s", w.name), w.damage, w.reach, w.tags))
}
