package world

import (
	"anvil/internal/core/definition"
	"anvil/internal/grid"
)

func (w *World) Width() int {
	return w.grid.Width()
}

func (w *World) Height() int {
	return w.grid.Height()
}

func (w *World) AddOccupant(pos grid.Position, o definition.Creature) {
	cell, _ := w.At(pos)
	cell.AddOccupant(o)
}

func (w *World) RemoveOccupant(pos grid.Position, o definition.Creature) {
	cell, _ := w.At(pos)
	cell.RemoveOccupant(o)
}

func (w *World) At(pos grid.Position) (definition.WorldCell, bool) {
	return w.grid.At(pos)
}

func (w World) IsValidPosition(pos grid.Position) bool {
	return w.grid.IsValidPosition(pos)
}

func (w World) CreaturesInRange(pos grid.Position, radius int) []definition.Creature {
	creatures := make([]definition.Creature, 0)
	for _, cell := range w.grid.CellsInRange(pos, radius) {
		creatures = append(creatures, cell.occupants...)
	}
	return creatures
}

func (w World) Navigation() definition.Pathfinding {
	return w.navigation
}
