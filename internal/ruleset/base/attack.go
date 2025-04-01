package base

import (
	"anvil/internal/core"
	"anvil/internal/tag"
)

type AttackAction struct {
	owner *core.Creature
}

func NewAttackAction(owner *core.Creature) AttackAction {
	return AttackAction{
		owner: owner,
	}
}

func (a AttackAction) Name() string {
	return "Attack"
}

func (a AttackAction) Perform(target *core.Creature) {
	if target == nil {
		return
	}
	a.owner.Log().Start(core.UseActionEventType, core.UseActionEvent{Action: a, Source: *a.owner, Target: *target})
	defer a.owner.Log().End()
	result := a.owner.AttackRoll(target, tag.NewContainer())
	if result.Success {
		target.TakeDamage(5)
	}
}
