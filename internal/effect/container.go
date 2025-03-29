package effect

import (
	"anvil/internal/effect/state"
)

type Container struct {
	effects []Effect
}

func NewContainer(effects ...Effect) *Container {
	return &Container{effects: effects}
}

func (c *Container) Add(effect Effect) {
	c.effects = append(c.effects, effect)
}

func (c *Container) Remove(effect Effect) {
	for i, e := range c.effects {
		if e.Name == effect.Name {
			c.effects = append(c.effects[:i], c.effects[i+1:]...)
			return
		}
	}
}

func (c *Container) Evaluate(state state.State) {
	for _, effect := range c.effects {
		effect.Evaluate(state)
	}
}
