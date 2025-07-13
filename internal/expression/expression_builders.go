package expression

import "anvil/internal/tag"

func FromConstant(value int, source string, components ...Component) *Expression {
	expr := &Expression{Rng: NewRngRoller()}
	expr.AddConstant(value, source, components...)
	return expr
}

func FromDice(times int, sides int, source string, components ...Component) *Expression {
	expr := &Expression{Rng: NewRngRoller()}
	expr.AddDice(times, sides, source, components...)
	return expr
}

func FromD20(source string, components ...Component) *Expression {
	expr := &Expression{Rng: NewRngRoller()}
	expr.AddD20(source, components...)
	return expr
}

func FromDamageConstant(value int, tags tag.Container, source string, components ...Component) *Expression {
	expr := &Expression{Rng: NewRngRoller()}
	expr.AddDamageConstant(value, tags, source, components...)
	return expr
}

func FromDamageDice(times int, sides int, tags tag.Container, source string, components ...Component) *Expression {
	expr := &Expression{Rng: NewRngRoller()}
	expr.AddDamageDice(times, sides, tags, source, components...)
	return expr
}

func (e *Expression) AddConstant(value int, source string, components ...Component) {
	e.Components = append(e.Components, newConstantComponent(value, tag.NewContainer(Primary), source, components...))
}

func (e *Expression) AddDice(times int, sides int, source string, components ...Component) {
	e.Components = append(e.Components, newDiceComponent(times, sides, tag.NewContainer(Primary), source, components...))
}

func (e *Expression) AddD20(source string, components ...Component) {
	e.Components = append(e.Components, newD20Component(tag.NewContainer(Primary), source, components...))
}

func (e *Expression) AddDamageConstant(value int, tags tag.Container, source string, components ...Component) {
	e.Components = append(e.Components, newConstantComponent(value, tags, source, components...))
}

func (e *Expression) AddDamageDice(times int, sides int, tags tag.Container, source string, components ...Component) {
	e.Components = append(e.Components, newDiceComponent(times, sides, tags, source, components...))
}

func FromDamageResult(damage Expression) *Expression {
	return damage.Clone()
}
