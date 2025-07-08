package expression

import (
	"slices"
)

func (e *Expression) EvaluateGroup() *Expression {
	result := Expression{Rng: DefaultRoller{}}
	e.Evaluate()
	groups := e.groupComponentsBy()
	for _, group := range groups {
		value := 0
		for _, component := range group {
			value += component.Value
		}
		result.AddDamageConstant(value, group[0].Source, group[0].Tags, group...)
	}

	if len(e.Components) > 0 && len(result.Components) > 0 {
		result.Components[0].IsCritical = e.Components[0].IsCritical
	}

	return result.Evaluate()
}

func (e *Expression) uniqueTags() []string {
	set := make([]string, 0, len(e.Components))
	for _, component := range e.Components {
		tags := e.primaryTags(component.Tags)
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
			tags := e.primaryTags(component.Tags)
			if tags.ID() != id {
				continue
			}

			components[i] = append(components[i], component)
		}
	}
	return components
}
