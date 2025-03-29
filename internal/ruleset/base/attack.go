package base

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/event"
	"anvil/internal/tagcontainer"
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
	defer a.owner.Log().End()
	result := a.owner.AttackRoll(target, tagcontainer.New())
	if result.Success {
		target.TakeDamage(5)
	}
}
