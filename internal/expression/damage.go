package expression

import (
	"fmt"
	"math"

	"anvil/internal/mathi"
	"anvil/internal/tag"
)

func (e *Expression) HalveDamage(tag tag.Tag, source string) {
	for i, component := range e.Components {
		if !component.Tags.MatchTag(tag) {
			continue
		}
		e.evaluateComponent(&component)
		value := math.Floor(float64(component.Value) / 2.0)
		e.Components[i] = Component{
			Type:       Constant,
			Source:     fmt.Sprintf("Halved (%s) %s", source, component.Source),
			Value:      int(value),
			Tags:       component.Tags.Clone(),
			Components: []Component{component},
		}
	}
}

func (e *Expression) ReplaceWith(value int, source string) {
	components := e.Components
	e.Components = []Component{{
		Type:       Constant,
		Source:     source,
		Value:      value,
		Tags:       tag.NewContainer(),
		Components: components,
	}}
}

func (e *Expression) DoubleDice(source string) {
	var components []Component
	for _, component := range e.Components {
		components = append(components, component)
		if !component.Type.Match(Dice) {
			continue
		}

		newComponent := component.Clone()
		newComponent.Source = source
		components = append(components, newComponent)
	}
	e.Components = components
}

func (e *Expression) MaxDice(source string) {
	var components []Component
	for _, component := range e.Components {
		components = append(components, component)
		if !component.Type.Match(Dice) {
			continue
		}

		newComponent := component.Clone()
		newComponent.Source = source
		newComponent.Type = Constant
		newComponent.Value = mathi.Abs(component.Sides * component.Times)
		components = append(components, newComponent)
	}
	e.Components = components
}

func (e *Expression) HasDamageType(t tag.Tag) bool {
	for _, component := range e.Components {
		if component.Tags.MatchTag(t) {
			return true
		}
	}

	return false
}
