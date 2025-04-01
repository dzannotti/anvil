package core

import (
	"anvil/internal/core/pathfinding"
	"anvil/internal/grid"
)

type World struct {
	Grid       *grid.Grid[WorldCell]
	Navigation *pathfinding.Pathfinding
}

func NewWorld(width int, height int) *World {
	return &World{
		Grid:       grid.New(width, height, NewWorldCell),
		Navigation: pathfinding.New(width, height),
	}
}

func (w *World) Width() int {
	return w.Grid.Width()
}

func (w *World) Height() int {
	return w.Grid.Height()
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
	return w.Grid.At(pos)
}

func (w World) IsValidPosition(pos grid.Position) bool {
	return w.Grid.IsValidPosition(pos)
}

func (w World) CreaturesInRange(pos grid.Position, radius int) []*Creature {
	creatures := make([]*Creature, 0)
	for _, cell := range w.Grid.CellsInRange(pos, radius) {
		creatures = append(creatures, cell.occupants...)
	}
	return creatures
}
