package expression

import "anvil/internal/tag"

func FromScalar(value int, source string, terms ...Term) Expression {
	e := Expression{rng: defaultRoller{}}
	e.AddScalar(value, source, terms...)
	return e
}

func FromDice(times int, sides int, source string, terms ...Term) Expression {
	e := Expression{rng: defaultRoller{}}
	e.AddDice(times, sides, source, terms...)
	return e
}

func FromD20(source string, terms ...Term) Expression {
	e := Expression{rng: defaultRoller{}}
	e.AddD20(source, terms...)
	return e
}

func FromDamageScalar(value int, source string, tags tag.Container, terms ...Term) Expression {
	e := Expression{rng: defaultRoller{}}
	e.AddDamageScalar(value, source, tags, terms...)
	return e
}

func FromDamageDice(times int, sides int, source string, tags tag.Container, terms ...Term) Expression {
	e := Expression{rng: defaultRoller{}}
	e.AddDamageDice(times, sides, source, tags, terms...)
	return e
}

func FromDamageResult(res Expression) Expression {
	e := Expression{rng: defaultRoller{}}
	e.Terms = append(e.Terms, res.Terms...)
	for i := range e.Terms {
		e.Terms[i].Terms = make([]Term, 0)
	}
	return e
}
