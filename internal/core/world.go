package core

import (
	"math"

	"anvil/internal/core/pathfinding"
	"anvil/internal/core/shapes"
	"anvil/internal/grid"
)

type World struct {
	Grid    *grid.Grid[WorldCell]
	Request *Request
}

func NewWorld(width int, height int) *World {
	return &World{
		Grid: grid.New(width, height, func(pos grid.Position) WorldCell {
			return WorldCell{Position: pos}
		}),
	}
}

func (w *World) Width() int {
	return w.Grid.Width
}

func (w *World) Height() int {
	return w.Grid.Height
}

func (w *World) AddOccupant(pos grid.Position, o *Actor) {
	cell, _ := w.At(pos)
	cell.AddOccupant(o)
}

func (w *World) RemoveOccupant(pos grid.Position, o *Actor) {
	cell, _ := w.At(pos)
	cell.RemoveOccupant(o)
}

func (w *World) At(pos grid.Position) (*WorldCell, bool) {
	return w.Grid.At(pos)
}

func (w World) IsValidPosition(pos grid.Position) bool {
	return w.Grid.IsValidPosition(pos)
}

func (w World) ActorsInRange(pos grid.Position, radius int, filter func(*Actor) bool) []*Actor {
	actors := make([]*Actor, 0, 10)
	cells := w.Grid.CellsInRange(pos, radius)
	for _, cell := range cells {
		other, ok := cell.Occupant()
		if !ok || !filter(other) {
			continue
		}
		actors = append(actors, other)
	}
	return actors
}

func (w World) ActorAt(pos grid.Position) (*Actor, bool) {
	cell, ok := w.At(pos)
	if !ok {
		return nil, false
	}
	return cell.Occupant()
}

func (w World) FindPath(start grid.Position, end grid.Position) (*pathfinding.Result, bool) {
	navCost := func(pos grid.Position) int {
		cell, _ := w.Grid.At(pos)
		if cell.Tile == Wall {
			return math.MaxInt
		}
		return 1
	}
	return pathfinding.FindPath(start, end, w.Width(), w.Height(), navCost)
}

func (w World) HasLineOfSight(from grid.Position, to grid.Position) bool {
	isDiagonalStep := func(a grid.Position, b grid.Position) bool {
		return a.X != b.X && a.Y != b.Y
	}

	isBlocked := func(pos grid.Position) bool {
		cell, ok := w.Grid.At(pos)
		if !ok {
			return true
		}
		return cell.Tile == Wall
	}

	line := shapes.Line(from, to)

	for i := 1; i < len(line); i++ {
		current := line[i]
		if isDiagonalStep(line[i-1], current) {
			adj1 := grid.Position{X: current.X, Y: line[i-1].Y}
			adj2 := grid.Position{X: line[i-1].X, Y: current.Y}
			if isBlocked(adj1) && isBlocked(adj2) {
				return false
			}
		}

		if isBlocked(current) {
			return false
		}
	}

	return true
}

func (w World) FloodFill(start grid.Position, radius int) []grid.Position {
	isBlocked := func(pos grid.Position) bool {
		cell, ok := w.Grid.At(pos)
		if !ok {
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
	if w.Request != nil {
		panic("There is already a pending request. Please wait until it is resolved.")
	}

	w.Request = &Request{
		Target:   actor,
		Text:     text,
		Options:  options,
		Response: make(chan RequestOption),
	}
	selectedOption := <-w.Request.Response
	w.Request = nil
	return selectedOption
}
