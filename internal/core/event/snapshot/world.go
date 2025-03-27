package snapshot

import (
	"anvil/internal/core/definition"
	"anvil/internal/grid"
)

type WorldCell struct {
	Walkable bool
	Occupant Creature
}

type World struct {
	Cells [][]WorldCell
}

func CaptureWorld(world definition.World) World {
	cells := make([][]WorldCell, world.Width())
	for x := 0; x < world.Width(); x++ {
		cells[x] = make([]WorldCell, world.Height())
		for y := 0; y < world.Height(); y++ {
			cell, _ := world.At(grid.NewPosition(x, y))
			path, _ := world.Navigation().At(grid.NewPosition(x, y))
			occupant, _ := cell.Occupant()
			cells[x][y] = WorldCell{Walkable: path.IsWalkable(), Occupant: CaptureCreature(occupant)}
		}
	}
	return World{
		Cells: cells,
	}
}
