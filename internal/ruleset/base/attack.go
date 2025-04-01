package base

import (
	"anvil/internal/core"
	"anvil/internal/tag"
)

type AttackAction struct {
	Owner *core.Creature
}

func NewAttackAction(owner *core.Creature) AttackAction {
	return AttackAction{
		Owner: owner,
	}
}

func (a AttackAction) Name() string {
	return "Attack"
}

func (a AttackAction) Perform(target *core.Creature) {
	if target == nil {
		return
	}
	a.Owner.Log.Start(core.UseActionEventType, core.UseActionEvent{Action: a, Source: *a.Owner, Target: *target})
	defer a.Owner.Log.End()
	result := a.Owner.AttackRoll(target, tag.NewContainer())
	if result.Success {
		target.TakeDamage(5)
	}
}
