package expression

import (
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
	if len(e.Components) == 0 {
		return inputTags
	}

	if inputTags.IsEmpty() || inputTags.HasTag(tag.FromString("primary")) {
		return e.Components[0].Tags
	}

	return inputTags
}
