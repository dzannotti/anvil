package expression

import (
	"anvil/internal/tag"
)

func (e *Expression) AddConstant(value int, source string, components ...Component) {
	e.Components = append(e.Components, Component{
		Type:       TypeConstant,
		Source:     source,
		Value:      value,
		Components: components,
	})
}

func (e *Expression) AddDice(times int, sides int, source string, components ...Component) {
	e.Components = append(e.Components, Component{
		Type:       TypeDice,
		Source:     source,
		Times:      times,
		Sides:      sides,
		Tags:       e.primaryTags(tag.NewContainerFromString("primary")),
		Components: components,
	})
}

func (e *Expression) AddD20(source string, components ...Component) {
	e.Components = append(e.Components, Component{
		Type:       TypeDice20,
		Source:     source,
		Times:      1,
		Sides:      20,
		Tags:       e.primaryTags(tag.NewContainerFromString("primary")),
		Components: components,
	})
}

func (e *Expression) AddDamageConstant(value int, source string, tags tag.Container, components ...Component) {
	e.Components = append(e.Components, Component{
		Type:       TypeDamageConstant,
		Source:     source,
		Value:      value,
		Tags:       e.primaryTags(tags),
		Components: components,
	})
}

func (e *Expression) AddDamageDice(times int, sides int, source string, tags tag.Container, components ...Component) {
	e.Components = append(e.Components, Component{
		Type:       TypeDamageDice,
		Source:     source,
		Times:      times,
		Sides:      sides,
		Tags:       e.primaryTags(tags),
		Components: components,
	})
}
