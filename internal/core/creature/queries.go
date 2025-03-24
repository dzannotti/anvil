package creature

func (c Creature) Name() string {
	return c.name
}

func (c Creature) IsDead() bool {
	return c.hitPoints == 0
}
