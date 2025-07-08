package expression

import (
	"fmt"

	"anvil/internal/mathi"
)

func (e *Expression) Evaluate() *Expression {
	e.Value = 0
	if e.Rng == nil {
		e.Rng = DefaultRoller{}
	}
	for i := range e.Components {
		e.evaluateComponent(&e.Components[i])
		e.Value += e.Components[i].Value
	}
	return e
}

func (e *Expression) evaluateComponent(component *Component) {
	if component.Type.Match(Constant) {
		return
	}
	e.evaluateDice(component)
}

func (e *Expression) evaluateDice(component *Component) {
	if !component.hasRollModifier() {
		e.evaluateDiceRoll(component)
		return
	}
	e.evaluateD20Roll(component)
}

func (e *Expression) evaluateDiceRoll(component *Component) {
	sign := mathi.Sign(component.Times)
	times := mathi.Abs(component.Times)
	component.Values = make([]int, times)
	component.Value = 0
	for i := range times {
		component.Values[i] = e.Rng.Roll(component.Sides)
		component.Value += component.Values[i]
	}
	component.Value *= sign
}

func (e *Expression) evaluateD20Roll(component *Component) {
	values := []int{e.Rng.Roll(component.Sides), e.Rng.Roll(component.Sides)}
	component.Values = values
	if len(component.HasAdvantage) > 0 {
		component.Value = mathi.Max(values[0], values[1])
		return
	}

	component.Value = mathi.Min(values[0], values[1])
}

func (e *Expression) IsCriticalSuccess() bool {
	if len(e.Components) == 0 {
		return false
	}

	first := e.Components[0]
	return first.IsCritical == CriticalSuccess || first.Value == first.Sides
}

func (e *Expression) IsCriticalFailure() bool {
	if len(e.Components) == 0 {
		return false
	}

	first := e.Components[0]
	if first.IsCritical == CriticalFailure {
		return true
	}

	if len(first.Values) == 0 {
		return false
	}

	return first.Values[0] == 1
}

func (e *Expression) SetCriticalSuccess(source string) {
	e.Components[0].IsCritical = CriticalSuccess
	e.Components[0].Source = fmt.Sprintf("%s as Critical success (%s)",
		e.Components[0].Source, source)
}

func (e *Expression) SetCriticalFailure(source string) {
	e.Components[0].IsCritical = CriticalFailure
	e.Components[0].Source = fmt.Sprintf("%s as Critical failure (%s)",
		e.Components[0].Source, source)
}
