package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"anvil/internal/core"
)

const cellSize = 64

type RenderCell struct {
	Tile core.TerrainType
}

type RenderWorld struct {
	width  int
	height int
	cells  []RenderCell
}

func NewRenderWorld(width int, height int) RenderWorld {
	return RenderWorld{
		width:  width,
		height: height,
		cells:  make([]RenderCell, width*height),
	}
}

func (r *RenderWorld) Render() {
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			r.RenderCell(x, y)
		}
	}
}

func (r *RenderWorld) RenderCell(x int, y int) {
	startX, startY := r.ToWorldPosition(x, y)

	rl.DrawRectangleLines(int32(startX), int32(startY), int32(cellSize), int32(cellSize), rl.Black)
}

func (r RenderWorld) ToWorldPosition(x int, y int) (int, int) {
	return x * cellSize, y * cellSize
}
