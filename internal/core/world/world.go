package world

import (
	"anvil/internal/grid"
)

type World struct {
	grid *grid.Grid[WorldCell]
}

func New(width int, height int) *World {
	return &World{
		grid: grid.New(width, height, NewWorldCell),
	}
}
