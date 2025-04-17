package ui

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

func DrawPoint(pos Vector2i, color Color, size int) {
	halfSize := size / 2
	FillRectangle(Rectangle{X: pos.X - halfSize, Y: pos.Y - halfSize, Width: size, Height: size}, color)
}

func DrawLine(start Vector2i, end Vector2i, color Color, thickness int) {
	rl.DrawLineEx(start.toRaylib(), end.toRaylib(), float32(thickness), color)
}
