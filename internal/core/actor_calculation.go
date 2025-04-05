package core

import (
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/expression"
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
		a.Log.Start(AttributeChangedType, AttributeChangeEvent{Source: a, Attribute: t, OldValue: old, Value: old + val, Reason: reason})
		defer a.Log.End()
		a.HitPoints = val
		a.Evaluate(AttributeChanged, AttributeChangedState{Source: a, Attribute: t, OldValue: old, Value: old + val})
		return
	}
	panic("ModifyAttribute not implemented")
}

func (a *Actor) SaveThrow(t tag.Tag, dc int) CheckResult {
	expr := expression.FromD20("Base")
	before := BeforeSavingThrowState{Expression: &expr, Source: a, Attribute: t, DifficultyClass: dc}
	a.Log.Start(SavingThrowType, SavingThrowEvent{Expression: &expr, Source: a, Attribute: t, DifficultyClass: dc})
	a.Evaluate(BeforeSavingThrow, &before)
	expr.Evaluate()
	after := AfterSavingThrowState{Result: &expr, Source: a, Attribute: t, DifficultyClass: dc}
	a.Evaluate(AfterSavingThrow, &after)
	ok := expr.Value >= dc
	crit := expr.IsCritical()
	a.Log.Add(ExpressionResultType, ExpressionResultEvent{Expression: expr})
	a.Log.Add(CheckResultType, CheckResultEvent{Value: expr.Value, Against: dc, Critical: crit, Success: ok})
	return CheckResult{Value: expr.Value, Against: dc, Critical: crit, Success: ok}
}
