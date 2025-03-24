package creature

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
