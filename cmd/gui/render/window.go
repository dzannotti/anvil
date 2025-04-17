package render

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Window struct{}

func (r Window) Open() {
	rl.InitWindow(1280, 720, "Anvil")
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
	rl.ClearBackground(rl.Black)
}

func (r Window) EndFrame() {
	rl.EndDrawing()
}
