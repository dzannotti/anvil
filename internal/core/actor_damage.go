package core

import (
	"anvil/internal/core/tags"
	"anvil/internal/expression"
	"anvil/internal/tag"
)

func (a *Actor) TakeDamage(damage int) {
	a.HitPoints = max(a.HitPoints-damage, 0)
	a.Log.Add(TakeDamageEventType, TakeDamageEvent{Target: *a, Damage: damage})
	if a.IsDead() {
		a.Log.Start(DeathEventType, DeathEvent{Actor: *a})
		a.Log.Add(ConfirmEventType, ConfirmEvent{Confirm: true})
		a.Log.End()
	}
}

func (a *Actor) AttackRoll(target *Actor, tc tag.Container) CheckResult {
	expr := expression.FromD20("Base")
	a.Log.Start(AttackRollEventType, AttackRollEvent{Source: *a, Target: *target})
	defer a.Log.End()
	before := BeforeAttackRollState{Source: a, Target: target, Expression: &expr, Tags: tc}
	a.Effects.Evaluate(BeforeAttackRollStateType, before)
	expr.Evaluate()
	after := AfterAttackRollState{Source: a, Target: target, Result: &expr, Tags: tc}
	a.Effects.Evaluate(AfterAttackRollStateType, after)
	a.Log.Add(ExpressionResultEventType, ExpressionResultEvent{Expression: expr})
	value := after.Result.Value
	crit := after.Result.IsCritical()
	targetAC := target.ArmorClass()
	a.Log.Add(AttributeCalculationEventType, AttributeCalculationEvent{Attribute: tags.ArmorClass, Expression: targetAC})
	ok := value >= targetAC.Value
	a.Log.Add(CheckResultEventType, CheckResultEvent{Value: value, Against: targetAC.Value, Critical: crit, Success: ok})
	return CheckResult{Value: value, Against: expr.Value, Critical: crit, Success: ok}
}
