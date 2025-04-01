package core

import (
	"anvil/internal/grid"
	"slices"
)

type WorldCell struct {
	Position  grid.Position
	Occupants []*Creature
}

func (c *WorldCell) AddOccupant(creature *Creature) {
	c.Occupants = append(c.Occupants, creature)
}

func (c *WorldCell) RemoveOccupant(creature *Creature) {
	c.Occupants = slices.DeleteFunc(c.Occupants, func(o *Creature) bool {
		return o == creature
	})
}

func (c *WorldCell) Occupant() (*Creature, bool) {
	if len(c.Occupants) == 0 {
		return &Creature{}, false
	}
	return c.Occupants[0], true
}

func (c *WorldCell) IsOccupied() bool {
	return len(c.Occupants) > 0
}
