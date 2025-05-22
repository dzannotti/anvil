package core

import "slices"

type EffectContainer struct {
	effects []*Effect
}

func (c *EffectContainer) Add(effect ...*Effect) {
	c.effects = append(c.effects, effect...)
	slices.SortFunc(c.effects, func(a, b *Effect) int {
		return int(a.Priority) - int(b.Priority)
	})
}

func (c *EffectContainer) Remove(effect *Effect) {
	for i, e := range c.effects {
		if e.Name == effect.Name {
			c.effects = append(c.effects[:i], c.effects[i+1:]...)
			return
		}
	}
}

func (c *EffectContainer) Evaluate(event string, state any) {
	for _, effect := range c.effects {
		effect.Evaluate(event, state)
	}
}

func (c *EffectContainer) Find(id string) *Effect {
	for _, effect := range c.effects {
		if effect.ID() == id {
			return effect
		}
	}
	return nil
}
