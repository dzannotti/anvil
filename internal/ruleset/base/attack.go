package base

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/event"
)

type AttackAction struct {
	owner definition.Creature
}

func NewAttackAction(owner definition.Creature) AttackAction {
	return AttackAction{
		owner: owner,
	}
}

func (a AttackAction) Name() string {
	return "Attack"
}

func (a AttackAction) Perform(target definition.Creature) {
	if target == nil {
		return
	}
	a.owner.Log().Start(event.NewUseAction(a, a.owner, target))
	target.TakeDamage(5)
	a.owner.Log().End()
}
