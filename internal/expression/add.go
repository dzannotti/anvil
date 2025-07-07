package expression

import (
	"anvil/internal/core/tags"
	"anvil/internal/tag"
)

func (e *Expression) AddConstant(value int, source string, components ...Component) {
	e.Components = append(e.Components, Component{
		Type:       TypeConstant,
		Source:     source,
		Value:      value,
		Tags:       tag.NewContainer(tags.ComponentConstant),
		Components: components,
	})
}

func (e *Expression) AddDice(times int, sides int, source string, components ...Component) {
	containerTags := tag.NewContainer(tags.ComponentDice)
	e.Components = append(e.Components, Component{
		Type:       TypeDice,
		Source:     source,
		Times:      times,
		Sides:      sides,
		Tags:       e.primaryTags(containerTags),
		Components: components,
	})
}

func (e *Expression) AddD20(source string, components ...Component) {
	containerTags := tag.NewContainer(tags.ComponentDice20)
	e.Components = append(e.Components, Component{
		Type:       TypeDice20,
		Source:     source,
		Times:      1,
		Sides:      20,
		Tags:       e.primaryTags(containerTags),
		Components: components,
	})
}

func (e *Expression) AddDamageConstant(value int, source string, componentTags tag.Container, components ...Component) {
	containerTags := componentTags.Clone()
	containerTags.AddTag(tags.ComponentDamageConstant)
	e.Components = append(e.Components, Component{
		Type:       TypeDamageConstant,
		Source:     source,
		Value:      value,
		Tags:       e.primaryTags(containerTags),
		Components: components,
	})
}

func (e *Expression) AddDamageDice(times int, sides int, source string, componentTags tag.Container, components ...Component) {
	containerTags := componentTags.Clone()
	containerTags.AddTag(tags.ComponentDamageDice)
	e.Components = append(e.Components, Component{
		Type:       TypeDamageDice,
		Source:     source,
		Times:      times,
		Sides:      sides,
		Tags:       e.primaryTags(containerTags),
		Components: components,
	})
}
