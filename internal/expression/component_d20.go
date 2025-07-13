package expression

import (
	"anvil/internal/mathi"
	"anvil/internal/tag"
)

type D20Component struct {
	value        int
	values       []int
	source       string
	advantage    []string
	disadvantage []string
	tags         tag.Container
	components   []Component
}

func newD20Component(tags tag.Container, source string, components ...Component) *D20Component {
	return &D20Component{tags: tags, source: source, components: components, advantage: []string{}, disadvantage: []string{}}
}

func (c *D20Component) Kind() ComponentKind {
	return ComponentKindD20
}

func (c *D20Component) Tags() tag.Container {
	return c.tags
}

func (c *D20Component) Value() int {
	return c.value
}

func (c *D20Component) Source() string {
	return c.source
}

func (c *D20Component) Components() []Component {
	return c.components
}

func (c *D20Component) Evaluate(ctx *Context) int {
	hasAdvantage := len(c.advantage) > 0
	hasDisadvantage := len(c.disadvantage) > 0
	isModified := (hasAdvantage || hasDisadvantage) && (!hasAdvantage || !hasDisadvantage)
	if !isModified {
		c.values = []int{ctx.Rng.Roll(20)}
		c.value = c.values[0]
		return c.value
	}

	values := []int{ctx.Rng.Roll(20), ctx.Rng.Roll(20)}
	c.values = values
	c.value = mathi.Min(values[0], values[1])

	if hasAdvantage {
		c.value = mathi.Max(values[0], values[1])
	}

	return c.value
}

func (c *D20Component) GiveAdvantage(source string) *D20Component {
	c.advantage = append(c.advantage, source)
	return c
}

func (c *D20Component) GiveDisadvantage(source string) *D20Component {
	c.disadvantage = append(c.disadvantage, source)
	return c
}

func (c *D20Component) IsCritical() bool {
	return c.IsCriticalSuccess() || c.IsCriticalFailure()
}

func (c *D20Component) IsCriticalSuccess() bool {
	return c.value == 20
}

func (c *D20Component) IsCriticalFailure() bool {
	return c.value == 1
}

func (c *D20Component) Advantage() []string {
	return c.advantage
}

func (c *D20Component) Disadvantage() []string {
	return c.disadvantage
}

func (c *D20Component) Expected() int {
	// Average of d20 is 10.5, round to 11
	return 11
}

func (c *D20Component) Clone() Component {
	components := make([]Component, len(c.components))
	for i, component := range c.components {
		components[i] = component.Clone()
	}
	return &D20Component{value: c.value, tags: c.tags.Clone(), source: c.source, components: components}
}
