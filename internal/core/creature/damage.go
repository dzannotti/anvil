package creature

func (c Creature) Attack(other *Creature) {
	other.TakeDamage(5)
}

func (c *Creature) TakeDamage(damage int) {
	c.hitPoints = max(c.hitPoints-damage, 0)
}

func (c *Creature) StartTurn() {

}
