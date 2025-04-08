package aiutils

import (
	"anvil/internal/core"
	"anvil/internal/grid"
	"fmt"
	"slices"
	"strings"
)

func debugValidPositions(a *core.Actor, valid []grid.Position) {
	sb := strings.Builder{}
	for y := range a.World.Height() {
		for x := 0; x < a.World.Width(); x++ {
			pos := grid.Position{X: x, Y: y}
			cell, _ := a.World.At(pos)
			if cell.Tile == core.Wall {
				sb.WriteString("#")
				continue
			}
			if slices.Contains(valid, pos) {
				if cell.IsOccupied() {
					sb.WriteString("~")
					continue
				}
				sb.WriteString("@")
				continue
			}
			if cell.IsOccupied() {
				occupant, _ := cell.Occupant()
				sb.WriteString(occupant.Name[0:1])
				continue
			}
			sb.WriteString(".")
		}
		sb.WriteString("\n")
	}
	fmt.Println(sb.String())
}
