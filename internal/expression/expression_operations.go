package expression

import (
	"anvil/internal/tag"
	"fmt"
)

func (e *Expression) GiveAdvantage(source string) {
	if len(e.Components) == 0 {
		panic("no components to give advantage to")
	}
	d20, ok := e.Components[0].(*D20Component)

	if !ok {
		panic("can only give advantage to d20 components")
	}

	d20.GiveAdvantage(source)
}

func (e *Expression) GiveDisadvantage(source string) {
	if len(e.Components) == 0 {
		panic("no components to give disadvantage to")
	}
	d20, ok := e.Components[0].(*D20Component)

	if !ok {
		panic("can only give disadvantage to d20 components")
	}

	d20.GiveDisadvantage(source)
}

func (e *Expression) ReplaceWith(value int, source string, tags tag.Container) {
	e.Components = []Component{newConstantComponent(value, tags, source)}
}

func (e *Expression) DoubleDice(source string) {
	for _, component := range e.Components {
		dice, ok := component.(*DiceComponent)
		if !ok {
			continue
		}

		diceSource := fmt.Sprintf("%s (%s)", dice.Source(), source)
		newDice := newDiceComponent(dice.Times(), dice.Sides(), dice.Tags(), diceSource, dice.Components()...)
		e.Components = append(e.Components, newDice)
	}
}

func (e *Expression) MaxDice(source string) {
	for _, component := range e.Components {
		dice, ok := component.(*DiceComponent)
		if !ok {
			continue
		}

		diceSource := fmt.Sprintf("%s (%s)", dice.Source(), source)
		newDice := newConstantComponent(dice.Times()*dice.Sides(), dice.Tags(), diceSource)
		e.Components = append(e.Components, newDice)
	}
}

func (e *Expression) EvaluateDamage() *Expression {
	e.Evaluate()

	if len(e.Components) == 0 {
		return e
	}

	groups := groupComponentsByTags(e.Components)
	if len(groups) == 0 {
		return e
	}

	newComponents := make([]Component, 0, len(groups))
	for _, group := range groups {
		groupTags := resolveGroupTags(group, e.Components)
		groupSource := buildGroupSource(group)
		groupValue := 0
		for _, comp := range group {
			groupValue += comp.Value()
		}

		newComponents = append(newComponents, newConstantComponent(groupValue, groupTags, groupSource))
	}

	e.Components = newComponents
	return e
}

func (e *Expression) IsCriticalSuccess() bool {
	d20, ok := e.Components[0].(*D20Component)
	if !ok {
		return false
	}

	return d20.IsCriticalSuccess()
}

func (e *Expression) IsCriticalFailure() bool {
	d20, ok := e.Components[0].(*D20Component)
	if !ok {
		return false
	}

	return d20.IsCriticalFailure()
}


func (e *Expression) Expected() int {
	total := 0
	for _, component := range e.Components {
		total += component.Expected()
	}
	return total
}

func (e *Expression) HasDamageType(damageType tag.Tag) bool {
	for _, component := range e.Components {
		compTags := component.Tags()
		if compTags.HasTag(damageType) {
			return true
		}
	}
	return false
}
