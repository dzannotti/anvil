package weapon

import (
	"fmt"

	"anvil/internal/core"
	"anvil/internal/ruleset/base"
)

type Weapon struct {
	name   string
	damage []core.DamageSource
}

func (w Weapon) Name() string {
	return w.name
}

func (w Weapon) OnEquip(a *core.Actor) {
	a.AddAction(base.NewAttackAction(a, fmt.Sprintf("Attack with %s", w.name), w.damage))
}
