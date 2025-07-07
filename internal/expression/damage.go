package expression

import (
	"fmt"
	"math"
	"strings"

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
		src := fmt.Sprintf("Halved (%s) %s", source, component.Source)
		e.Components[i] = Component{
			Type:       TypeConstant,
			Source:     src,
			Value:      int(value),
			Tags:       component.Tags,
			Components: []Component{component},
		}
	}
}

func (e *Expression) ReplaceWith(value int, source string) {
	components := e.Components
	e.Components = []Component{{
		Type:       TypeConstant,
		Source:     source,
		Value:      value,
		Components: components,
	}}
}

func (e *Expression) DoubleDice(source string) {
	components := []Component{}
	for _, component := range e.Components {
		components = append(components, component)
		if !strings.Contains(string(component.Type), string(TypeDice)) {
			continue
		}
		newComponent := component.Clone()
		newComponent.Source = source
		components = append(components, newComponent)
	}
	e.Components = components
}

func (e *Expression) MaxDice(source string) {
	components := []Component{}
	for _, component := range e.Components {
		components = append(components, component)
		if !strings.Contains(string(component.Type), string(TypeDice)) {
			continue
		}
		newComponent := component.Clone()
		newComponent.Source = source
		newComponent.Type = TypeConstant
		newComponent.Value = mathi.Abs(component.Sides * component.Times)
		components = append(components, newComponent)
	}
	e.Components = components
}

func (e *Expression) IsDamageType(t tag.Tag) bool {
	for _, component := range e.Components {
		if component.Tags.MatchTag(t) {
			return true
		}
	}
	return false
}
