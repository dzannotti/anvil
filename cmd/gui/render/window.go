package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Window struct {
	Width  int
	Height int
}

func (r Window) Open() {
	r.Width = 1280
	r.Height = 720
	rl.InitWindow(int32(r.Width), int32(r.Height), "Anvil")
	rl.SetTargetFPS(60)
}

func (r Window) Close() {
	defer rl.CloseWindow()
}

func (r Window) ShouldClose() bool {
	return rl.WindowShouldClose()
}

func (r Window) StartFrame() {
	rl.BeginDrawing()
	rl.ClearBackground(Base)
}

func (r Window) EndFrame() {
	rl.EndDrawing()
}
