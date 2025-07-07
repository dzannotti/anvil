package expression

import (
	"anvil/internal/core/tags"
	"anvil/internal/tag"
)

type Expression struct {
	Components []Component
	Value      int
	Rng        DiceRoller
}

func (e *Expression) Clone() Expression {
	components := make([]Component, len(e.Components))
	for i := range e.Components {
		components[i] = e.Components[i].Clone()
	}
	return Expression{
		Value:      e.Value,
		Components: components,
		Rng:        e.Rng,
	}
}

func (e *Expression) primaryTags(inputTags tag.Container) tag.Container {
	if len(e.Components) > 0 {
		if inputTags.IsEmpty() || inputTags.HasTag(tag.FromString("primary")) {
			return e.Components[0].Tags
		}
	}
	return inputTags
}

func (e *Expression) primaryTagsForGrouping(inputTags tag.Container) tag.Container {
	var resultTags tag.Container
	
	if len(e.Components) > 0 {
		if inputTags.IsEmpty() || inputTags.HasTag(tag.FromString("primary")) {
			resultTags = e.Components[0].Tags.Clone()
		} else {
			resultTags = inputTags.Clone()
		}
	} else {
		resultTags = inputTags.Clone()
	}
	
	// Filter out component type tags for grouping purposes
	resultTags.RemoveTag(tags.ComponentType, tags.ComponentConstant, tags.ComponentDamageConstant, tags.ComponentDice, tags.ComponentDice20, tags.ComponentDamageDice)
	return resultTags
}
