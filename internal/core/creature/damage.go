package creature

import (
	"anvil/internal/core/event"
)

func (c *Creature) TakeDamage(damage int) {
	c.hitPoints = max(c.hitPoints-damage, 0)
	c.log.Add(event.NewTakeDamage(c, damage))
}

func (c *Creature) StartTurn() {

}
