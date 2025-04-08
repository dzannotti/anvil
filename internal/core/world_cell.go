package core

import (
	"slices"

	"anvil/internal/grid"
)

type WorldCell struct {
	Position  grid.Position
	Occupants []*Actor
}

func (c *WorldCell) AddOccupant(actor *Actor) {
	c.Occupants = append(c.Occupants, actor)
}

func (c *WorldCell) RemoveOccupant(actor *Actor) {
	c.Occupants = slices.DeleteFunc(c.Occupants, func(o *Actor) bool {
		return o == actor
	})
}

func (c *WorldCell) Occupant() (*Actor, bool) {
	if len(c.Occupants) == 0 {
		return nil, false
	}
	return c.Occupants[0], true
}

func (c *WorldCell) IsOccupied() bool {
	return len(c.Occupants) > 0
}
