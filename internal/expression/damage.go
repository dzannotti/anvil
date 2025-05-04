package expression

import (
	"fmt"
	"math"
	"strings"

	"anvil/internal/mathi"
	"anvil/internal/tag"
)

func (e *Expression) HalveDamage(tag tag.Tag, source string) {
	for i, term := range e.Terms {
		if !term.Tags.MatchTag(tag) {
			continue
		}
		e.evaluateTerm(&term)
		value := math.Floor(float64(term.Value) / 2.0)
		src := fmt.Sprintf("Halved (%s) %s", source, term.Source)
		e.Terms[i] = makeTerm(TypeScalarHalve, src, term)
		e.Terms[i].Value = int(value)
		e.Terms[i].Tags = term.Tags
	}
}

func (e *Expression) ReplaceWith(value int, source string) {
	terms := e.Terms
	newTerm := makeTerm(TypeScalarReplace, source, terms...)
	newTerm.Value = value
	e.Terms = []Term{newTerm}
}

func (e *Expression) DoubleDice(source string) {
	terms := []Term{}
	for _, term := range e.Terms {
		terms = append(terms, term)
		if !strings.Contains(string(term.Type), string(TypeDice)) {
			continue
		}
		newTerm := term.Clone()
		newTerm.Source = source
		terms = append(terms, newTerm)
	}
	e.Terms = terms
}

func (e *Expression) MaxDice(source string) {
	terms := []Term{}
	for _, term := range e.Terms {
		terms = append(terms, term)
		if !strings.Contains(string(term.Type), string(TypeDice)) {
			continue
		}
		newTerm := term.Clone()
		newTerm.Source = source
		newTerm.Type = TypeScalarMax
		newTerm.Value = mathi.Abs(term.Sides * term.Times)
		terms = append(terms, newTerm)
	}
	e.Terms = terms
}

func (e Expression) IsDamageType(t tag.Tag) bool {
	for _, term := range e.Terms {
		if term.Tags.MatchTag(t) {
			return true
		}
	}
	return false
}
