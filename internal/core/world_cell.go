package core

import (
	"anvil/internal/grid"
	"slices"
)

type WorldCell struct {
	position  grid.Position
	occupants []*Creature
}

func NewWorldCell(position grid.Position) WorldCell {
	return WorldCell{
		position:  position,
		occupants: make([]*Creature, 0),
	}
}

func (c *WorldCell) Position() grid.Position {
	return c.position
}

func (c *WorldCell) AddOccupant(creature *Creature) {
	c.occupants = append(c.occupants, creature)
}

func (c *WorldCell) RemoveOccupant(creature *Creature) {
	c.occupants = slices.DeleteFunc(c.occupants, func(o *Creature) bool {
		return o == creature
	})
}

func (c *WorldCell) Occupant() (*Creature, bool) {
	if len(c.occupants) == 0 {
		return &Creature{}, false
	}
	return c.occupants[0], true
}

func (c *WorldCell) IsOccupied() bool {
	return len(c.occupants) > 0
}
