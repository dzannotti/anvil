package expression

import (
	"fmt"
	"strings"

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
	if strings.Contains(string(component.Type), string(TypeConstant)) {
		return
	}
	e.evaluateDice(component)
}

func (e *Expression) evaluateDice(component *Component) {
	if !component.shouldModifyRoll() {
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
	return e.Components[0].IsCritical == 1 || e.Components[0].Value == e.Components[0].Sides
}

func (e *Expression) IsCriticalFailure() bool {
	if len(e.Components) == 0 {
		return false
	}
	return e.Components[0].IsCritical == -1 || e.Components[0].Values[0] == 1
}

func (e *Expression) SetCriticalSuccess(source string) {
	e.Components[0].IsCritical = 1
	e.Components[0].Source += fmt.Sprintf(" as Critical success (%s)", source)
}

func (e *Expression) SetCriticalFailure(source string) {
	e.Components[0].IsCritical = -1
	e.Components[0].Source += fmt.Sprintf(" as Critical failure (%s)", source)
}
