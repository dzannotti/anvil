package world

import (
	"anvil/internal/core/creature"
	"anvil/internal/core/definition"
	"anvil/internal/grid"
	"slices"
)

type Cell struct {
	position  grid.Position
	occupants []definition.Creature
}

func NewCell(position grid.Position) Cell {
	return Cell{
		position:  position,
		occupants: make([]definition.Creature, 0),
	}
}

func (c *Cell) Position() grid.Position {
	return c.position
}

func (c *Cell) AddOccupant(creature definition.Creature) {
	c.occupants = append(c.occupants, creature)
}

func (c *Cell) RemoveOccupant(creature definition.Creature) {
	c.occupants = slices.DeleteFunc(c.occupants, func(o definition.Creature) bool {
		return o == creature
	})
}

func (c *Cell) Occupant() (definition.Creature, bool) {
	if len(c.occupants) == 0 {
		return &creature.Creature{}, false
	}
	return c.occupants[0], true
}

func (c *Cell) IsOccupied() bool {
	return len(c.occupants) > 0
}
