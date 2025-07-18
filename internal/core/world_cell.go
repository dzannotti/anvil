package core

import (
	"slices"

	"anvil/internal/grid"
)

type TerrainType int

const (
	Normal TerrainType = iota
	Wall
)

type WorldCell struct {
	Position  grid.Position
	Tile      TerrainType
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

func (c *WorldCell) Occupant() *Actor {
	if len(c.Occupants) == 0 {
		return nil
	}
	return c.Occupants[0]
}

func (c *WorldCell) IsOccupied() bool {
	return len(c.Occupants) > 0
}
