package world

import (
	"anvil/internal/core/pathfinding"
	"anvil/internal/grid"
)

type World struct {
	grid       *grid.Grid[Cell]
	navigation *pathfinding.Pathfinding
}

func New(width int, height int) *World {
	return &World{
		grid:       grid.New(width, height, NewCell),
		navigation: pathfinding.New(width, height),
	}
}
