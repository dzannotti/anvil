package core

import (
	"math"

	"anvil/internal/core/pathfinding"
	"anvil/internal/core/shapes"
	"anvil/internal/grid"
	"anvil/internal/loader"
)

type World struct {
	Grid            *grid.Grid[WorldCell]
	lineOfSightCalc *LineOfSightCalculator
	requestManager  *RequestManager
}

func NewWorld(definition loader.WorldDefinition) *World {
	w := &World{
		Grid: grid.New(definition.Width, definition.Height, func(pos grid.Position) WorldCell {
			return WorldCell{Position: pos}
		}),
		requestManager: NewRequestManager(),
	}
	w.lineOfSightCalc = NewLineOfSightCalculator(w)
	return w
}

func (w *World) Width() int {
	return w.Grid.Width
}

func (w *World) Height() int {
	return w.Grid.Height
}

func (w *World) AddOccupant(pos grid.Position, o *Actor) {
	w.At(pos).AddOccupant(o)
}

func (w *World) RemoveOccupant(pos grid.Position, o *Actor) {
	w.At(pos).RemoveOccupant(o)
}

func (w *World) At(pos grid.Position) *WorldCell {
	return w.Grid.At(pos)
}

func (w *World) IsValidPosition(pos grid.Position) bool {
	return w.Grid.IsValidPosition(pos)
}

func (w *World) ActorsInRange(pos grid.Position, radius int, filter func(*Actor) bool) []*Actor {
	actors := make([]*Actor, 0, 10)
	cells := w.Grid.Cells(shapes.Square(pos, radius, radius))
	for _, cell := range cells {
		other := cell.Occupant()
		if other == nil || !filter(other) {
			continue
		}

		actors = append(actors, other)
	}
	return actors
}

func (w *World) ActorAt(pos grid.Position) *Actor {
	if !w.IsValidPosition(pos) {
		return nil
	}

	return w.At(pos).Occupant()
}

func (w *World) FindPath(start grid.Position, end grid.Position) (*pathfinding.Result, bool) {
	navCost := func(pos grid.Position) int {
		cell := w.Grid.At(pos)
		if cell.Tile == Wall {
			return math.MaxInt
		}

		return 1
	}
	result := pathfinding.FindPath(start, end, w.Width(), w.Height(), navCost)
	return result, result.Found
}

func (w *World) HasLineOfSight(from grid.Position, to grid.Position) bool {
	return w.lineOfSightCalc.HasLineOfSight(from, to)
}

func (w *World) FloodFill(start grid.Position, radius int) []grid.Position {
	isBlocked := func(pos grid.Position) bool {
		cell := w.Grid.At(pos)
		if cell == nil {
			return true
		}

		if cell.Tile == Wall {
			return true
		}

		return false
	}
	return shapes.FloodFill(start, radius, isBlocked)
}

func (w *World) Ask(actor *Actor, text string, options []RequestOption) RequestOption {
	result, err := w.requestManager.Ask(actor, text, options)
	if err != nil {
		panic(err.Error())
	}

	return result
}

func (w *World) RequestManager() *RequestManager {
	return w.requestManager
}
