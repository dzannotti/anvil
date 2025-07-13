package expression

import "anvil/internal/tag"

type ConstantComponent struct {
	value      int
	source     string
	tags       tag.Container
	components []Component
}

func newConstantComponent(value int, tags tag.Container, source string, components ...Component) *ConstantComponent {
	return &ConstantComponent{value: value, tags: tags, source: source, components: components}
}

func (c *ConstantComponent) Kind() ComponentKind {
	return ComponentKindConstant
}

func (c *ConstantComponent) Tags() tag.Container {
	return c.tags
}

func (c *ConstantComponent) Value() int {
	return c.value
}

func (c *ConstantComponent) Source() string {
	return c.source
}

func (c *ConstantComponent) Components() []Component {
	return c.components
}

func (c *ConstantComponent) Evaluate(_ *Context) int {
	return c.value
}

func (c *ConstantComponent) Expected() int {
	return c.value
}

func (c *ConstantComponent) Clone() Component {
	components := make([]Component, len(c.components))
	for i, component := range c.components {
		components[i] = component.Clone()
	}
	return &ConstantComponent{value: c.value, tags: c.tags.Clone(), source: c.source, components: components}
}
