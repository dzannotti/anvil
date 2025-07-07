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
	w.At(pos).AddOccupant(o)
}

func (w *World) RemoveOccupant(pos grid.Position, o *Actor) {
	w.At(pos).RemoveOccupant(o)
}

func (w *World) At(pos grid.Position) *WorldCell {
	return w.Grid.At(pos)
}

func (w World) IsValidPosition(pos grid.Position) bool {
	return w.Grid.IsValidPosition(pos)
}

func (w World) ActorsInRange(pos grid.Position, radius int, filter func(*Actor) bool) []*Actor {
	actors := make([]*Actor, 0, 10)
	cells := w.Grid.CellsInRange(pos, radius)
	for _, cell := range cells {
		other := cell.Occupant()
		if other == nil || !filter(other) {
			continue
		}
		actors = append(actors, other)
	}
	return actors
}

func (w World) ActorAt(pos grid.Position) *Actor {
	if !w.IsValidPosition(pos) {
		return nil
	}
	return w.At(pos).Occupant()
}

func (w World) FindPath(start grid.Position, end grid.Position) (*pathfinding.Result, bool) {
	navCost := func(pos grid.Position) int {
		cell := w.Grid.At(pos)
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
		cell := w.Grid.At(pos)
		if cell == nil {
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
