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
		Expression: &expr,
		Attribute:  tags.ArmorClass,
	}
	a.Evaluate(AttributeCalculation, &s)
	s.Expression.Evaluate()
	return s.Expression
}

func (a *Actor) Attribute(t tag.Tag) *expression.Expression {
	expr := expression.FromScalar(a.Attributes.Value(t), tags.ToReadableTag(t))
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
