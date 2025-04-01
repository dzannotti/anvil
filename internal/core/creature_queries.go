package core

import (
	"anvil/internal/eventbus"
)

func (c Creature) Log() *eventbus.Hub {
	return c.log
}

func (c Creature) Name() string {
	return c.name
}

func (c Creature) IsDead() bool {
	return c.hitPoints == 0
}

func (c Creature) HitPoints() int {
	return c.hitPoints
}

func (c Creature) MaxHitPoints() int {
	return c.maxHitPoints
}

func (c Creature) Actions() []Action {
	return c.actions
}

func (c *Creature) AddAction(action ...Action) {
	c.actions = append(c.actions, action...)
}
