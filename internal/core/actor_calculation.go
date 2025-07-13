package core

import (
	"math"

	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/expression"
	"anvil/internal/grid"
	"anvil/internal/mathi"
	"anvil/internal/tag"
)

func (a *Actor) ArmorClass() *expression.Expression {
	expr := expression.FromConstant(10, "Base")
	dex := a.Attribute(tags.AttributeDexterity)
	expr.AddConstant(stats.AttributeModifier(dex.Value), "Attribute Modifier", dex.Components...)
	s := AttributeCalculation{
		Source:     a,
		Expression: expr,
		Attribute:  tags.ActorArmorClass,
	}
	a.Evaluate(&s)
	s.Expression.Evaluate()
	return s.Expression
}

func (a *Actor) Attribute(t tag.Tag) *expression.Expression {
	expr := expression.FromConstant(a.Attributes.Value(t), tags.ToReadable(t))
	s := AttributeCalculation{
		Expression: expr,
		Attribute:  t,
	}
	a.Evaluate(&s)
	s.Expression.Evaluate()
	return s.Expression
}

func (a *Actor) Proficiency(tags tag.Container) int {
	return a.Proficiencies.Value(tags)
}

func (a *Actor) ModifyAttribute(t tag.Tag, val int, reason string) {
	if t.MatchExact(tags.ActorHitPoints) {
		old := a.HitPoints
		a.Dispatcher.Begin(AttributeChangeEvent{Source: a, Attribute: t, OldValue: old, Value: old + val, Reason: reason})
		defer a.Dispatcher.End()
		a.HitPoints = val
		a.Evaluate(&AttributeChanged{Source: a, Attribute: t, OldValue: old, Value: old + val})
		return
	}

	panic("ModifyAttribute not implemented")
}

func (a *Actor) SaveThrow(t tag.Tag, dc int) CheckResult {
	expr := expression.FromD20("Base")
	before := PreSavingThrow{Expression: expr, Source: a, Attribute: t, DifficultyClass: dc}
	a.Dispatcher.Begin(SavingThrowEvent{Expression: expr, Source: a, Attribute: t, DifficultyClass: dc})
	defer a.Dispatcher.End()
	a.Evaluate(&before)
	expr.Evaluate()
	after := PostSavingThrow{Result: expr, Source: a, Attribute: t, DifficultyClass: dc}
	a.Evaluate(&after)
	success := expr.Value >= dc
	crit := false
	if after.Result.IsCriticalSuccess() {
		crit = true
		success = true
	}
	if after.Result.IsCriticalFailure() {
		crit = true
	}
	a.Dispatcher.Emit(ExpressionResultEvent{Expression: expr})
	a.Dispatcher.Emit(SavingThrowResultEvent{Actor: a, Value: expr.Value, Against: dc, Critical: crit, Success: success})
	return CheckResult{Value: expr.Value, Against: dc, Critical: crit, Success: success}
}

func (a *Actor) TakeDamage(damage expression.Expression) {
	expr := expression.FromDamageResult(damage)
	before := PreTakeDamage{Expression: expr, Source: a}
	a.Evaluate(&before)
	res := expr.Evaluate()
	actual := a.HitPoints - mathi.Clamp(a.HitPoints-res.Value, 0, math.MaxInt)
	a.HitPoints = mathi.Clamp(a.HitPoints-actual, 0, math.MaxInt)
	a.Dispatcher.Begin(TakeDamageEvent{Target: a, Damage: expr})
	after := PostTakeDamage{Result: res, Source: a, ActualDamage: actual}
	a.Effects.Evaluate(&after)
	a.Dispatcher.End()
}

func (a *Actor) AttackRoll(target *Actor, tc tag.Container) CheckResult {
	expr := expression.FromD20("Base")
	a.Dispatcher.Begin(AttackRollEvent{Source: a, Target: target})
	defer a.Dispatcher.End()
	before := PreAttackRoll{Source: a, Target: target, Expression: expr, Tags: tc}
	a.Effects.Evaluate(&before)
	expr.Evaluate()
	after := PostAttackRoll{Source: a, Target: target, Result: expr, Tags: tc}
	a.Effects.Evaluate(&after)
	a.Dispatcher.Emit(ExpressionResultEvent{Expression: expr})
	value := after.Result.Value
	targetAC := target.ArmorClass()
	a.Dispatcher.Emit(AttributeCalculationEvent{Attribute: tags.ActorArmorClass, Expression: targetAC})
	hit := value >= targetAC.Value
	crit := false
	if after.Result.IsCriticalSuccess() {
		crit = true
		hit = true
	}
	if after.Result.IsCriticalFailure() {
		crit = true
		hit = false
	}
	a.Dispatcher.Emit(CheckResultEvent{Actor: a, Value: value, Against: targetAC.Value, Critical: crit, Success: hit, Tags: tc})
	return CheckResult{Value: value, Against: targetAC.Value, Critical: crit, Success: hit}
}

func (a *Actor) DamageRoll(ds DamageSource, crit bool) *expression.Expression {
	expr := ds.Damage().Clone()
	a.Dispatcher.Begin(DamageRollEvent{Source: a, DamageSource: ds})
	defer a.Dispatcher.End()
	before := PreDamageRoll{Source: a, Expression: expr, Tags: *ds.Tags(), Critical: crit}
	a.Effects.Evaluate(&before)
	res := expr.EvaluateDamage()
	a.Dispatcher.Emit(ExpressionResultEvent{Expression: res})
	after := PostDamageRoll{Source: a, Result: res, Tags: *ds.Tags(), Critical: crit}
	a.Effects.Evaluate(&after)
	return res
}

func (a *Actor) Move(to grid.Position, action Action) {
	a.Dispatcher.Begin(MoveStepEvent{World: a.World, Source: a, From: a.Position, To: to})
	defer a.Dispatcher.End()
	before := PreMoveStep{
		Source:  a,
		From:    a.Position,
		To:      to,
		CanMove: true,
		Action:  action,
	}
	a.Effects.Evaluate(&before)
	a.Dispatcher.Emit(ConfirmEvent{Actor: a, Confirm: before.CanMove})
	if before.CanMove {
		a.World.RemoveOccupant(a.Position, a)
		a.Position = to
		a.World.AddOccupant(to, a)
	}
}
