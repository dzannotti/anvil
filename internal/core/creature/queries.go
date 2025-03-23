package creature

import (
	"anvil/internal/core/team"
)

func (c Creature) Name() string {
	return c.name
}

func (c Creature) HitPoints() int {
	return c.hitPoints
}

func (c Creature) MaxHitPoints() int {
	return c.maxHitPoints
}

func (c Creature) Team() team.Team {
	return c.team
}

func (c Creature) ActionPoints() int {
	return c.actionPoints
}

func (c Creature) IsDead() bool {
	return c.hitPoints == 0
}
