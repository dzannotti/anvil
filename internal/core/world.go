package core

import (
	"anvil/internal/core/definition"
	"anvil/internal/core/pathfinding"
	"anvil/internal/grid"
)

type World struct {
	grid       *grid.Grid[WorldCell]
	navigation *pathfinding.Pathfinding
}

func NewWorld(width int, height int) *World {
	return &World{
		grid:       grid.New(width, height, NewWorldCell),
		navigation: pathfinding.New(width, height),
	}
}

func (w *World) Width() int {
	return w.grid.Width()
}

func (w *World) Height() int {
	return w.grid.Height()
}

func (w *World) AddOccupant(pos grid.Position, o *Creature) {
	cell, _ := w.At(pos)
	cell.AddOccupant(o)
}

func (w *World) RemoveOccupant(pos grid.Position, o *Creature) {
	cell, _ := w.At(pos)
	cell.RemoveOccupant(o)
}

func (w *World) At(pos grid.Position) (*WorldCell, bool) {
	return w.grid.At(pos)
}

func (w World) IsValidPosition(pos grid.Position) bool {
	return w.grid.IsValidPosition(pos)
}

func (w World) CreaturesInRange(pos grid.Position, radius int) []*Creature {
	creatures := make([]*Creature, 0)
	for _, cell := range w.grid.CellsInRange(pos, radius) {
		creatures = append(creatures, cell.occupants...)
	}
	return creatures
}

func (w World) Navigation() definition.Pathfinding {
	return w.navigation
}
