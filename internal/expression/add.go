package expression

import (
	"anvil/internal/tag"
)

func (e *Expression) AddScalar(value int, source string, terms ...Term) {
	term := makeTerm(TypeScalar, source, terms...)
	term.Value = value
	e.Terms = append(e.Terms, term)
}

func (e *Expression) AddDice(times int, sides int, source string, terms ...Term) {
	term := makeTerm(TypeDice, source, terms...)
	term.Times = times
	term.Sides = sides
	e.Terms = append(e.Terms, term)
}

func (e *Expression) AddD20(source string, terms ...Term) {
	term := makeTerm(TypeDice20, source, terms...)
	term.Times = 1
	term.Sides = 20
	e.Terms = append(e.Terms, term)
}

func (e *Expression) AddDamageScalar(value int, source string, tags tag.Container, terms ...Term) {
	term := makeTerm(TypeDamageScalar, source, terms...)
	term.Value = value
	term.Tags = e.primaryTags(tags)
	e.Terms = append(e.Terms, term)
}

func (e *Expression) AddDamageDice(times int, sides int, source string, tags tag.Container, terms ...Term) {
	term := makeTerm(TypeDamageDice, source, terms...)
	term.Times = times
	term.Sides = sides
	term.Tags = e.primaryTags(tags)
	e.Terms = append(e.Terms, term)
}

func (e Expression) primaryTags(tags tag.Container) tag.Container {
	if len(e.Terms) > 0 {
		if e.Terms[0].Tags.HasTag(tag.FromString("primary")) {
			return e.Terms[0].Tags
		}
	}
	return tags
}
