package expression

import (
	"strings"

	"github.com/adam-lavrik/go-imath/ix"
)

func (e *Expression) Evaluate() *Expression {
	e.Value = 0
	if e.rng == nil {
		e.rng = defaultRoller{}
	}
	for i := range e.Terms {
		e.evaluateTerm(&e.Terms[i])
		e.Value += e.Terms[i].Value
	}
	return e
}

func (e Expression) evaluateTerm(term *Term) {
	if strings.Contains(string(term.Type), string(TypeScalar)) {
		return
	}
	e.evaluateDice(term)
}

func (e Expression) evaluateDice(term *Term) {
	if !term.shouldModifyRoll() {
		e.evaluateDiceRoll(term)
		return
	}
	e.evaluateD20Roll(term)
}

func (e Expression) evaluateDiceRoll(term *Term) {
	sign := ix.Sign(term.Times)
	times := ix.Abs(term.Times)
	term.Values = make([]int, times)
	term.Value = 0
	for i := 0; i < times; i++ {
		term.Values[i] = e.rng.Roll(term.Sides)
		term.Value += term.Values[i]
	}
	term.Value *= sign
}

func (e Expression) evaluateD20Roll(term *Term) {
	values := []int{e.rng.Roll(term.Sides), e.rng.Roll(term.Sides)}
	term.Values = values
	if len(term.HasAdvantage) > 0 {
		term.Value = max(values[0], values[1])
		return
	}
	term.Value = min(values[0], values[1])
}

func (e Expression) IsCritical() bool {
	if len(e.Terms) == 0 {
		return false
	}
	if e.Terms[0].IsCritical != 0 {
		return e.Terms[0].IsCritical == 1
	}
	return e.Terms[0].Value == e.Terms[0].Sides
}
