package core

import (
	"github.com/adam-lavrik/go-imath/ix"

	"anvil/internal/core/tags"
	"anvil/internal/expression"
	"anvil/internal/tag"
)

func (a *Actor) Die() {
	a.Log.Start(DeathType, DeathEvent{Actor: *a})
	defer a.Log.End()
	a.Log.Add(ConfirmType, ConfirmEvent{Confirm: true})
}

func (a *Actor) TakeDamage(damage int) {
	expr := expression.FromDamageScalar(damage, "FIXME", tag.ContainerFromTag(tags.NaturalWeapon))
	crit := false
	before := BeforeTakeDamageState{Expression: &expr, Source: a, Critical: &crit}
	a.Evaluate(BeforeTakeDamage, &before)
	res := before.Expression.Evaluate()
	effective := a.HitPoints - ix.Max(a.HitPoints-res.Value, 0)
	a.HitPoints = ix.Max(a.HitPoints-effective, 0)
	a.Log.Start(TakeDamageType, TakeDamageEvent{Target: *a, Damage: damage})
	after := AfterTakeDamageState{Result: res, Source: a, Critical: &crit}
	a.Effects.Evaluate(AfterTakeDamage, &after)
	a.Log.End()
}

func (a *Actor) AttackRoll(target *Actor, tc tag.Container) CheckResult {
	expr := expression.FromD20("Base")
	a.Log.Start(AttackRollType, AttackRollEvent{Source: *a, Target: *target})
	defer a.Log.End()
	before := BeforeAttackRollState{Source: a, Target: target, Expression: &expr, Tags: tc}
	a.Effects.Evaluate(BeforeAttackRoll, &before)
	expr.Evaluate()
	after := AfterAttackRollState{Source: a, Target: target, Result: &expr, Tags: tc}
	a.Effects.Evaluate(AfterAttackRoll, &after)
	a.Log.Add(ExpressionResultType, ExpressionResultEvent{Expression: expr})
	value := after.Result.Value
	crit := after.Result.IsCritical()
	targetAC := target.ArmorClass()
	a.Log.Add(AttributeCalculationType, AttributeCalculationEvent{Attribute: tags.ArmorClass, Expression: targetAC})
	ok := value >= targetAC.Value
	a.Log.Add(CheckResultType, CheckResultEvent{Value: value, Against: targetAC.Value, Critical: crit, Success: ok})
	return CheckResult{Value: value, Against: expr.Value, Critical: crit, Success: ok}
}

func (a *Actor) DamageRoll(ds []DamageSource, crit bool) *expression.Expression {
	expr := expression.Expression{}
	for _, d := range ds {
		expr.AddDamageDice(d.Times, d.Sides, d.Source, d.Tags)
	}
	a.Log.Start(DamageRollType, DamageRollEvent{Source: *a, DamageSource: ds})
	defer a.Log.End()
	before := BeforeDamageRollState{Source: a, Expression: &expr, Critical: &crit}
	a.Effects.Evaluate(BeforeDamageRoll, &before)
	expr.EvaluateGroup()
	a.Log.Add(ExpressionResultType, ExpressionResultEvent{Expression: expr})
	after := AfterDamageRollState{Source: a, Result: &expr, Critical: &crit}
	a.Effects.Evaluate(AfterDamageRoll, &after)

	return &expr
}
