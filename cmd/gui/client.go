package main

import (
	"fmt"
	"net"

	"anvil/cmd/gui/render"
	"anvil/cmd/gui/ui"
)

func client(conn net.Conn) {
	window := render.Window{}
	window.Open()
	defer window.Close()
	ui.Init()
	defer ui.Close()

	for !window.ShouldClose() {
		window.StartFrame()

		ui.FillRectangle(ui.Rectangle{X: 1, Y: 1, Width: 100, Height: 100}, ui.Black)
		ui.DrawPoint(ui.Vector2i{X: 100, Y: 100}, ui.Red, 20)
		ui.DrawLine(ui.Vector2i{X: 100, Y: 100}, ui.Vector2i{X: 200, Y: 200}, ui.Blue, 5)
		//world.Render()
		//ui.DrawText("Congrats! You created your first window!", ui.Vector2i{X: 190, Y: 200}, ui.Black, 20)
		ui.FillRectangle(ui.Rectangle{X: 100, Y: 100, Width: 200, Height: 100}, ui.Green)
		ui.DrawString("Hello, World!", ui.Rectangle{X: 100, Y: 100, Width: 200, Height: 100}, ui.Black, 20, ui.AlignBottom)
		ui.DrawButton(ui.Rectangle{X: 300, Y: 100, Width: 200, Height: 100}, "My Button", ui.AlignBottom, 20, func() {
			fmt.Println("Clicked!")
		}, true)
		window.EndFrame()
		ui.ProcessInput()
	}
}
