package expression

import (
	"anvil/internal/tag"
)

func (e *Expression) AddConstant(value int, source string, components ...Component) {
	e.Components = append(e.Components, Component{
		Type:       Constant,
		Source:     source,
		Value:      value,
		Tags:       tag.NewContainer(),
		Components: components,
	})
}

func (e *Expression) AddDice(times int, sides int, source string, components ...Component) {
	e.Components = append(e.Components, Component{
		Type:       Dice,
		Source:     source,
		Times:      times,
		Sides:      sides,
		Tags:       e.primaryTags(tag.NewContainer()),
		Components: components,
	})
}

func (e *Expression) AddD20(source string, components ...Component) {
	e.Components = append(e.Components, Component{
		Type:       D20,
		Source:     source,
		Times:      1,
		Sides:      20,
		Tags:       e.primaryTags(tag.NewContainer()),
		Components: components,
	})
}

func (e *Expression) AddDamageConstant(value int, source string, tags tag.Container, components ...Component) {
	e.Components = append(e.Components, Component{
		Type:       DamageConstant,
		Source:     source,
		Value:      value,
		Tags:       e.primaryTags(tags),
		Components: components,
	})
}

func (e *Expression) AddDamageDice(times int, sides int, source string, tags tag.Container, components ...Component) {
	e.Components = append(e.Components, Component{
		Type:       DamageDice,
		Source:     source,
		Times:      times,
		Sides:      sides,
		Tags:       e.primaryTags(tags),
		Components: components,
	})
}

func (e *Expression) addD20Modifier(source string, modifiers *[]string) {
	if len(e.Components) == 0 {
		panic("cannot modify expression with no components")
	}

	if e.Components[0].Type != D20 {
		panic("cannot modify expression with non-d20 component")
	}

	*modifiers = append(*modifiers, source)
}

func (e *Expression) WithAdvantage(source string) {
	e.addD20Modifier(source, &e.Components[0].HasAdvantage)
}

func (e *Expression) WithDisadvantage(source string) {
	e.addD20Modifier(source, &e.Components[0].HasDisadvantage)
}
