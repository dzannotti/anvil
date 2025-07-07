package expression

import "anvil/internal/tag"

func FromConstant(value int, source string, components ...Component) Expression {
	e := Expression{Rng: DefaultRoller{}}
	e.AddConstant(value, source, components...)
	return e
}

func FromDice(times int, sides int, source string, components ...Component) Expression {
	e := Expression{Rng: DefaultRoller{}}
	e.AddDice(times, sides, source, components...)
	return e
}

func FromD20(source string, components ...Component) Expression {
	e := Expression{Rng: DefaultRoller{}}
	e.AddD20(source, components...)
	return e
}

func FromDamageConstant(value int, source string, tags tag.Container, components ...Component) Expression {
	e := Expression{Rng: DefaultRoller{}}
	e.AddDamageConstant(value, source, tags, components...)
	return e
}

func FromDamageDice(times int, sides int, source string, tags tag.Container, components ...Component) Expression {
	e := Expression{Rng: DefaultRoller{}}
	e.AddDamageDice(times, sides, source, tags, components...)
	return e
}

func FromDamageResult(res Expression) Expression {
	return res.Clone()
}
