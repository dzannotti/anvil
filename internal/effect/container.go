package effect

import "slices"

type Container struct {
	effects []Effect
}

func (c *Container) Add(effect Effect) {
	c.effects = append(c.effects, effect)
	slices.SortFunc(c.effects, func(a, b Effect) int {
		return int(a.Priority) - int(b.Priority)
	})
}

func (c *Container) Remove(effect Effect) {
	for i, e := range c.effects {
		if e.Name == effect.Name {
			c.effects = append(c.effects[:i], c.effects[i+1:]...)
			return
		}
	}
}

func (c *Container) Evaluate(event string, state any) {
	for _, effect := range c.effects {
		effect.Evaluate(event, state)
	}
}
