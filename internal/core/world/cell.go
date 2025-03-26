package world

import (
	"anvil/internal/core/definition"
	"anvil/internal/grid"
	"errors"
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

func (c *WorldCell) Occupant() (definition.Creature, error) {
	if len(c.occupants) == 0 {
		return nil, errors.New("no occupants")
	}
	return c.occupants[0], nil
}

func (c *WorldCell) IsOccupied() bool {
	return len(c.occupants) > 0
}
