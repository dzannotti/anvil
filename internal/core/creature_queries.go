package core

func (c Creature) IsDead() bool {
	return c.HitPoints == 0
}

func (c *Creature) AddAction(action ...Action) {
	c.Actions = append(c.Actions, action...)
}
