package render

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func FillRectangle(rect Rectangle, color Color) {
	rl.DrawRectangle(int32(rect.X), int32(rect.Y), int32(rect.Width), int32(rect.Height), color)
}

func FillRectangleGradient(rect Rectangle, color1 Color, color2 Color) {
	rl.DrawRectangleGradientV(int32(rect.X), int32(rect.Y), int32(rect.Width), int32(rect.Height), color1, color2)
}

func DrawRectangle(rect Rectangle, color Color, thickness int) {
	rl.DrawRectangleLinesEx(rect.toRaylib(), float32(thickness), color)
}

func FillCircle(pos Vector2i, radius int, color Color) {
	rl.DrawCircle(int32(pos.X), int32(pos.Y), float32(radius)*0.5, color)
}

func DrawLine(start Vector2i, end Vector2i, color Color, thickness int) {
	rl.DrawLineEx(start.ToRaylib(), end.ToRaylib(), float32(thickness), color)
}
