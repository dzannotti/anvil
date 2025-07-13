package expression

import (
	"anvil/internal/mathi"
	"anvil/internal/tag"
)

type DiceComponent struct {
	value      int
	values     []int
	source     string
	times      int
	sides      int
	tags       tag.Container
	components []Component
}

func newDiceComponent(times int, sides int, tags tag.Container, source string, components ...Component) *DiceComponent {
	return &DiceComponent{times: times, sides: sides, tags: tags, source: source, components: components}
}

func (c *DiceComponent) Kind() ComponentKind {
	return ComponentKindDice
}

func (c *DiceComponent) Tags() tag.Container {
	return c.tags
}

func (c *DiceComponent) Value() int {
	return c.value
}

func (c *DiceComponent) Source() string {
	return c.source
}

func (c *DiceComponent) Times() int {
	return c.times
}

func (c *DiceComponent) Sides() int {
	return c.sides
}

func (c *DiceComponent) Components() []Component {
	return c.components
}

func (c *DiceComponent) Values() []int {
	return c.values
}

func (c *DiceComponent) Evaluate(ctx *Context) int {
	c.values = []int{}
	times := mathi.Abs(c.times)
	for i := 0; i < times; i++ {
		c.values = append(c.values, ctx.Rng.Roll(c.sides))
	}

	c.value = mathi.Sum(c.values...) * mathi.Sign(c.times)
	return c.value
}

func (c *DiceComponent) Expected() int {
	// Average of a die is (sides + 1) / 2
	// For times dice: times * (sides + 1) / 2
	absTimesDie := mathi.Abs(c.times)
	expected := absTimesDie * (c.sides + 1) / 2
	return expected * mathi.Sign(c.times)
}

func (c *DiceComponent) Clone() Component {
	components := make([]Component, len(c.components))
	for i, component := range c.components {
		components[i] = component.Clone()
	}
	return &DiceComponent{times: c.times, sides: c.sides, tags: c.tags.Clone(), source: c.source, components: components}
}
