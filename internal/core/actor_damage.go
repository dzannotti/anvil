package core

import (
	"anvil/internal/core/tags"
	"anvil/internal/expression"
	"anvil/internal/tag"

	"github.com/adam-lavrik/go-imath/ix"
)

func (a *Actor) Die() {
	a.Log.Start(DeathEventType, DeathEvent{Actor: *a})
	defer a.Log.End()
	a.Log.Add(ConfirmEventType, ConfirmEvent{Confirm: true})
}

func (a *Actor) TakeDamage(damage int) {
	expr := expression.FromDamageScalar(damage, "FIXME", tag.ContainerFromTag(tags.ItemWeaponNatural))
	crit := false
	before := BeforeTakeDamageState{Expression: &expr, Source: a, Critical: &crit}
	a.Evaluate(BeforeTakeDamage, &before)
	res := before.Expression.Evaluate()
	effective := a.HitPoints - ix.Max(a.HitPoints-res.Value, 0)
	a.HitPoints = ix.Max(a.HitPoints-effective, 0)
	a.Log.Start(TakeDamageEventType, TakeDamageEvent{Target: *a, Damage: damage})
	after := AfterTakeDamageState{Result: res, Source: a, Critical: &crit}
	a.Effects.Evaluate(AfterTakeDamage, &after)
	a.Log.End()
}

func (a *Actor) AttackRoll(target *Actor, tc tag.Container) CheckResult {
	expr := expression.FromD20("Base")
	a.Log.Start(AttackRollEventType, AttackRollEvent{Source: *a, Target: *target})
	defer a.Log.End()
	before := BeforeAttackRollState{Source: a, Target: target, Expression: &expr, Tags: tc}
	a.Effects.Evaluate(BeforeAttackRoll, before)
	expr.Evaluate()
	after := AfterAttackRollState{Source: a, Target: target, Result: &expr, Tags: tc}
	a.Effects.Evaluate(AfterAttackRoll, after)
	a.Log.Add(ExpressionResultEventType, ExpressionResultEvent{Expression: expr})
	value := after.Result.Value
	crit := after.Result.IsCritical()
	targetAC := target.ArmorClass()
	a.Log.Add(AttributeCalculationEventType, AttributeCalculationEvent{Attribute: tags.ArmorClass, Expression: targetAC})
	ok := value >= targetAC.Value
	a.Log.Add(CheckResultEventType, CheckResultEvent{Value: value, Against: targetAC.Value, Critical: crit, Success: ok})
	return CheckResult{Value: value, Against: expr.Value, Critical: crit, Success: ok}
}
