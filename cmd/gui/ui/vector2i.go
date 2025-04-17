package ui

import rl "github.com/gen2brain/raylib-go/raylib"

type Vector2i struct {
	X int
	Y int
}

func (v Vector2i) toRaylib() rl.Vector2 {
	return rl.NewVector2(float32(v.X), float32(v.Y))
}
