package render

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"anvil/internal/grid"
)

type Camera struct {
	Camera    rl.Camera2D
	MoveSpeed float32
	ZoomSpeed float32
	MinZoom   float32
	IsEnabled bool
}

const (
	MoveSpeed = 5.0
	ZoomSpeed = 0.1
	Zoom      = 0.1
)

func (c *Camera) Update() {
	if rl.IsKeyDown(rl.KeyW) {
		c.Camera.Target.Y -= MoveSpeed / c.Camera.Zoom
	}
	if rl.IsKeyDown(rl.KeyS) {
		c.Camera.Target.Y += MoveSpeed / c.Camera.Zoom
	}
	if rl.IsKeyDown(rl.KeyA) {
		c.Camera.Target.X -= MoveSpeed / c.Camera.Zoom
	}
	if rl.IsKeyDown(rl.KeyD) {
		c.Camera.Target.X += MoveSpeed / c.Camera.Zoom
	}

	wheelMove := rl.GetMouseWheelMove()
	if wheelMove != 0 {
		mouseWorldPos := rl.GetScreenToWorld2D(rl.GetMousePosition(), c.Camera)
		c.Camera.Zoom += wheelMove * c.ZoomSpeed * c.Camera.Zoom

		if c.Camera.Zoom < c.MinZoom {
			c.Camera.Zoom = c.MinZoom
		}

		if wheelMove != 0 {
			mousePos := rl.GetMousePosition()
			newMouseWorldPos := rl.GetScreenToWorld2D(mousePos, c.Camera)
			c.Camera.Target.X += mouseWorldPos.X - newMouseWorldPos.X
			c.Camera.Target.Y += mouseWorldPos.Y - newMouseWorldPos.Y
		}
	}
}

func (c *Camera) Begin() {
	rl.BeginMode2D(c.Camera)
}

func (c *Camera) End() {
	rl.EndMode2D()
}

func (c *Camera) Reset(screenWidth int, screenHeight int) {
	c.Camera.Target = rl.Vector2{X: 0, Y: 0}
	c.Camera.Offset = rl.Vector2{X: float32(screenWidth) / 2, Y: float32(screenHeight) / 2}
	c.Camera.Rotation = 0.0
	c.Camera.Zoom = 1.0
}

func (c *Camera) GetWorldMousePosition() rl.Vector2 {
	return rl.GetScreenToWorld2D(rl.GetMousePosition(), c.Camera)
}

func (c *Camera) GetWorldScreenPos(pos Vector2i) rl.Vector2 {
	return rl.GetWorldToScreen2D(pos.ToRaylib(), c.Camera)
}

func (c *Camera) GetMouseGridPosition() grid.Position {
	mouseWorldPos := c.GetWorldMousePosition()
	return grid.Position{X: int(mouseWorldPos.X / CellSize), Y: int(mouseWorldPos.Y / CellSize)}
}

func (c *Camera) SetPosition(x, y float32) {
	c.Camera.Target = rl.Vector2{X: x, Y: y}
}

func (c *Camera) SetZoom(zoom float32) {
	c.Camera.Zoom = zoom
	if c.Camera.Zoom < c.MinZoom {
		c.Camera.Zoom = c.MinZoom
	}
}
