package world

import (
	"anvil/internal/core/creature"
	"anvil/internal/core/definition"
	"anvil/internal/grid"
	"slices"
)

type WorldCell struct {
	position  grid.Position
	occupants []definition.Creature
}

func NewWorldCell(position grid.Position) WorldCell {
	return WorldCell{
		position:  position,
		occupants: make([]definition.Creature, 0),
	}
}

func (c *WorldCell) Position() grid.Position {
	return c.position
}

func (c *WorldCell) AddOccupant(creature definition.Creature) {
	c.occupants = append(c.occupants, creature)
}

func (c *WorldCell) RemoveOccupant(creature definition.Creature) {
	slices.DeleteFunc(c.occupants, func(o definition.Creature) bool {
		return o == creature
	})
}

func (c *WorldCell) Occupant() (definition.Creature, bool) {
	if len(c.occupants) == 0 {
		return &creature.Creature{}, false
	}
	return c.occupants[0], true
}

func (c *WorldCell) IsOccupied() bool {
	return len(c.occupants) > 0
}
