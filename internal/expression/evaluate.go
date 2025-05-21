package expression

import (
	"fmt"
	"strings"

	"anvil/internal/mathi"
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
	sign := mathi.Sign(term.Times)
	times := mathi.Abs(term.Times)
	term.Values = make([]int, times)
	term.Value = 0
	for i := range times {
		term.Values[i] = e.rng.Roll(term.Sides)
		term.Value += term.Values[i]
	}
	term.Value *= sign
}

func (e Expression) evaluateD20Roll(term *Term) {
	values := []int{e.rng.Roll(term.Sides), e.rng.Roll(term.Sides)}
	term.Values = values
	if len(term.HasAdvantage) > 0 {
		term.Value = mathi.Max(values[0], values[1])
		return
	}
	term.Value = mathi.Min(values[0], values[1])
}

func (e Expression) IsCriticalSuccess() bool {
	if len(e.Terms) == 0 {
		return false
	}
	return e.Terms[0].IsCritical == 1 || e.Terms[0].Value == e.Terms[0].Sides
}

func (e Expression) IsCriticalFailure() bool {
	if len(e.Terms) == 0 {
		return false
	}
	return e.Terms[0].IsCritical == -1 || e.Terms[0].Values[0] == 1
}

func (e Expression) SetCriticalSuccess(source string) {
	e.Terms[0].IsCritical = 1
	e.Terms[0].Source = e.Terms[0].Source + fmt.Sprintf(" as Critical success (%s)", source)
}

func (e Expression) SetCriticalFailure(source string) {
	e.Terms[0].IsCritical = -1
	e.Terms[0].Source = e.Terms[0].Source + fmt.Sprintf(" as Critical failure (%s)", source)
}
