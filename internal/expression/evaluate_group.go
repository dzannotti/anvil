package expression

import (
	"slices"
	
	"anvil/internal/core/tags"
)

func (e *Expression) EvaluateGroup() *Expression {
	out := Expression{Rng: DefaultRoller{}}
	e.Evaluate()
	groups := e.groupComponentsBy()
	for _, group := range groups {
		value := 0
		for _, component := range group {
			value += component.Value
		}
		// Filter out component type tags to avoid double-adding them  
		filteredTags := group[0].Tags.Clone()
		filteredTags.RemoveTag(tags.ComponentType, tags.ComponentConstant, tags.ComponentDamageConstant, tags.ComponentDice, tags.ComponentDice20, tags.ComponentDamageDice)
		out.AddDamageConstant(value, group[0].Source, filteredTags, group...)
	}
	out.Components[0].IsCritical = e.Components[0].IsCritical
	return out.Evaluate()
}

func (e *Expression) uniqueTags() []string {
	set := make([]string, 0, len(e.Components))
	for _, component := range e.Components {
		tags := e.primaryTagsForGrouping(component.Tags)
		if slices.Contains(set, tags.ID()) {
			continue
		}
		set = append(set, tags.ID())
	}
	return set
}

func (e *Expression) groupComponentsBy() [][]Component {
	ids := e.uniqueTags()
	components := make([][]Component, len(ids))
	for i, id := range ids {
		for _, component := range e.Components {
			tags := e.primaryTagsForGrouping(component.Tags)
			if tags.ID() != id {
				continue
			}
			components[i] = append(components[i], component)
		}
	}
	return components
}
