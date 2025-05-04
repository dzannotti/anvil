package ui

import "anvil/internal/core"

func DrawRequest(world *core.World, window Window) {
	if world.Request == nil {
		return
	}
	FillRectangle(Rectangle{X: 0, Y: 0, Width: window.Width, Height: window.Height}, Color{R: Overlay0.R, G: Overlay0.G, B: Overlay0.B, A: 128})
}

