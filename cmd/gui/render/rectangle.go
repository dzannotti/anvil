package ui

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Rectangle struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (r Rectangle) toRaylib() rl.Rectangle {
	return rl.NewRectangle(float32(r.X), float32(r.Y), float32(r.Width), float32(r.Height))
}

func (r Rectangle) Position() Vector2i {
	return Vector2i{X: r.X, Y: r.Y}
}

func (r Rectangle) IsMouseOver() bool {
	mousePos := rl.GetMousePosition()
	return mousePos.X >= float32(r.X) && mousePos.X <= float32(r.X+r.Width) && mousePos.Y >= float32(r.Y) &&
		mousePos.Y <= float32(r.Y+r.Height)
}

func (r Rectangle) Expand(x int, y int) Rectangle {
	return Rectangle{X: r.X - x, Y: r.Y - y, Width: r.Width + x*2, Height: r.Height + y*2}
}
