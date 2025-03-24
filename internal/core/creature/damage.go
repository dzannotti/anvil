package creature

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/event"
)

func (c *Creature) Attack(other definition.Creature) {
	c.log.Start(event.NewUseAction("attack", c, other))
	other.TakeDamage(5)
	c.log.End()
}

func (c *Creature) TakeDamage(damage int) {
	c.hitPoints = max(c.hitPoints-damage, 0)
	c.log.Add(event.NewTakeDamage(c, damage))
}

func (c *Creature) StartTurn() {

}
