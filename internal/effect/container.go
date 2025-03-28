package effect

import (
	"anvil/internal/effect/state"
	"sync"
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

func (c *Container) Evaluate(state state.State, wg *sync.WaitGroup) {
	for _, effect := range c.effects {
		lwg := &sync.WaitGroup{}
		lwg.Add(1)
		effect.Evaluate(state, lwg)
		lwg.Wait()
	}
	wg.Done()
}
