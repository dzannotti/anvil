package core

import (
	"anvil/internal/core/shapes"
	"anvil/internal/grid"
)

type LineOfSightCalculator struct {
	world *World
}

func NewLineOfSightCalculator(world *World) *LineOfSightCalculator {
	return &LineOfSightCalculator{world: world}
}

func (los *LineOfSightCalculator) HasLineOfSight(from grid.Position, to grid.Position) bool {
	isDiagonalStep := func(a grid.Position, b grid.Position) bool {
		return a.X != b.X && a.Y != b.Y
	}

	isBlocked := func(pos grid.Position) bool {
		cell := los.world.Grid.At(pos)
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
