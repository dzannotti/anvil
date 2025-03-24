package team

import (
	"slices"

	"anvil/internal/core/creature"
)

func (t Team) IsDead() bool {
	for _, c := range t.members {
		if !c.IsDead() {
			return false
		}
	}
	return true
}

func (t Team) Name() string {
	return t.name
}

func (t Team) Contains(c *creature.Creature) bool {
	return slices.Contains(t.members, c)
}
