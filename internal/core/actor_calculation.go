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
	expr := expression.FromScalar(10, "Base")
	dex := a.Attribute(tags.Dexterity)
	expr.AddScalar(stats.AttributeModifier(dex.Value), "Attribute Modifier", dex.Terms...)
	s := AttributeCalculationState{
		Source:     a,
		Expression: &expr,
		Attribute:  tags.ArmorClass,
	}
	a.Evaluate(AttributeCalculation, &s)
	s.Expression.Evaluate()
	return s.Expression
}

func (a *Actor) Attribute(t tag.Tag) *expression.Expression {
	expr := expression.FromScalar(a.Attributes.Value(t), tags.ToReadable(t))
	s := AttributeCalculationState{
		Expression: &expr,
		Attribute:  t,
	}
	a.Evaluate(AttributeCalculation, &s)
	s.Expression.Evaluate()
	return s.Expression
}

func (a *Actor) Proficiency(tags tag.Container) int {
	return a.Proficiencies.Value(tags)
}

func (a *Actor) ModifyAttribute(t tag.Tag, val int, reason string) {
	if t.MatchExact(tags.HitPoints) {
		old := a.HitPoints
		a.Log.Start(
			AttributeChangedType,
			AttributeChangeEvent{Source: a, Attribute: t, OldValue: old, Value: old + val, Reason: reason},
		)
		defer a.Log.End()
		a.HitPoints = val
		a.Evaluate(AttributeChanged, &AttributeChangedState{Source: a, Attribute: t, OldValue: old, Value: old + val})
		return
	}
	panic("ModifyAttribute not implemented")
}

func (a *Actor) SaveThrow(t tag.Tag, dc int) CheckResult {
	expr := expression.FromD20("Base")
	before := BeforeSavingThrowState{Expression: &expr, Source: a, Attribute: t, DifficultyClass: dc}
	a.Log.Start(SavingThrowType, SavingThrowEvent{Expression: &expr, Source: a, Attribute: t, DifficultyClass: dc})
	defer a.Log.End()
	a.Evaluate(BeforeSavingThrow, &before)
	expr.Evaluate()
	after := AfterSavingThrowState{Result: &expr, Source: a, Attribute: t, DifficultyClass: dc}
	a.Evaluate(AfterSavingThrow, &after)
	ok := expr.Value >= dc
	crit := false
	if after.Result.IsCriticalSuccess() {
		crit = true
		ok = true
	}
	if after.Result.IsCriticalFailure() {
		crit = true
	}
	a.Log.Add(ExpressionResultType, ExpressionResultEvent{Expression: &expr})
	a.Log.Add(
		SavingThrowResultType,
		SavingThrowResultEvent{Actor: a, Value: expr.Value, Against: dc, Critical: crit, Success: ok},
	)
	return CheckResult{Value: expr.Value, Against: dc, Critical: crit, Success: ok}
}

func (a *Actor) TakeDamage(damage expression.Expression) {
	expr := expression.FromDamageResult(damage)
	before := BeforeTakeDamageState{Expression: &expr, Source: a}
	a.Evaluate(BeforeTakeDamage, &before)
	res := expr.Evaluate()
	actual := a.HitPoints - mathi.Clamp(a.HitPoints-res.Value, 0, math.MaxInt)
	a.HitPoints = mathi.Clamp(a.HitPoints-actual, 0, math.MaxInt)
	a.Log.Start(TakeDamageType, TakeDamageEvent{Target: a, Damage: &expr})
	after := AfterTakeDamageState{Result: res, Source: a, ActualDamage: actual}
	a.Effects.Evaluate(AfterTakeDamage, &after)
	a.Log.End()
}

func (a *Actor) AttackRoll(target *Actor, tc tag.Container) CheckResult {
	expr := expression.FromD20("Base")
	a.Log.Start(AttackRollType, AttackRollEvent{Source: a, Target: target})
	defer a.Log.End()
	before := BeforeAttackRollState{Source: a, Target: target, Expression: &expr, Tags: tc}
	a.Effects.Evaluate(BeforeAttackRoll, &before)
	expr.Evaluate()
	after := AfterAttackRollState{Source: a, Target: target, Result: &expr, Tags: tc}
	a.Effects.Evaluate(AfterAttackRoll, &after)
	a.Log.Add(ExpressionResultType, ExpressionResultEvent{Expression: &expr})
	value := after.Result.Value
	targetAC := target.ArmorClass()
	a.Log.Add(AttributeCalculationType, AttributeCalculationEvent{Attribute: tags.ArmorClass, Expression: targetAC})
	ok := value >= targetAC.Value
	crit := false
	if after.Result.IsCriticalSuccess() {
		crit = true
		ok = true
	}
	if after.Result.IsCriticalFailure() {
		crit = true
		ok = false
	}
	a.Log.Add(
		CheckResultType,
		CheckResultEvent{Actor: a, Value: value, Against: targetAC.Value, Critical: crit, Success: ok, Tags: tc},
	)
	return CheckResult{Value: value, Against: targetAC.Value, Critical: crit, Success: ok}
}

func (a *Actor) DamageRoll(ds []DamageSource, crit bool) *expression.Expression {
	expr := expression.Expression{}
	for _, d := range ds {
		expr.AddDamageDice(d.Times, d.Sides, d.Source, d.Tags)
	}
	if crit {
		expr.SetCriticalSuccess("Attack Roll")
	}
	a.Log.Start(DamageRollType, DamageRollEvent{Source: a, DamageSource: ds})
	defer a.Log.End()
	before := BeforeDamageRollState{Source: a, Expression: &expr}
	a.Effects.Evaluate(BeforeDamageRoll, &before)
	res := expr.EvaluateGroup()
	a.Log.Add(ExpressionResultType, ExpressionResultEvent{Expression: res})
	after := AfterDamageRollState{Source: a, Result: res}
	a.Effects.Evaluate(AfterDamageRoll, &after)
	return res
}

func (a *Actor) Move(to grid.Position, action Action) {
	a.Log.Start(MoveStepType, MoveStepEvent{World: a.World, Source: a, From: a.Position, To: to})
	defer a.Log.End()
	before := MoveState{
		Source:  a,
		From:    a.Position,
		To:      to,
		CanMove: true,
		Action:  action,
	}
	a.Effects.Evaluate(BeforeMoveStep, &before)
	a.Log.Add(ConfirmType, ConfirmEvent{Actor: a, Confirm: before.CanMove})
	if before.CanMove {
		a.World.RemoveOccupant(a.Position, a)
		a.Position = to
		a.World.AddOccupant(to, a)
	}
}
