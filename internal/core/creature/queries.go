package creature

import (
	"anvil/internal/core/definition"
	"anvil/internal/log"
)

func (c Creature) Log() *log.EventLog {
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

func (c Creature) Actions() []definition.Action {
	return c.actions
}

func (c *Creature) AddAction(action ...definition.Action) {
	c.actions = append(c.actions, action...)
}
