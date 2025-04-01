package base

import (
	"anvil/internal/core"
	"anvil/internal/core/definition"
	"anvil/internal/tag"
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
	a.owner.Log().Start(core.UseActionEventType, core.UseActionEvent{Action: a, Source: a.owner, Target: target})
	defer a.owner.Log().End()
	result := a.owner.AttackRoll(target, tag.NewContainer())
	if result.Success {
		target.TakeDamage(5)
	}
}
