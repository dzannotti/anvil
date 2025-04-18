package main

import (
	"fmt"
	"net"

	"anvil/cmd/gui/render"
	"anvil/cmd/gui/ui"
	"anvil/internal/demo"
	"anvil/internal/eventbus"
)

func client(_ net.Conn) {
	hub := eventbus.Hub{}
	hub.Subscribe(func(msg eventbus.Message) {})

	window := render.Window{}
	window.Open()
	defer window.Close()
	ui.Init()
	defer ui.Close()
	world, encounter := demo.Create(&hub)
	camera := render.Camera{}
	camera.Reset(window.Width, window.Height)
	camera.SetPosition(-20, -20)
	am := render.ActionManager{}

	endTurn := func() {
		encounter.EndTurn()
		am.SetActive(nil)
	}

	encounter.Start()
	for !window.ShouldClose() {
		window.StartFrame()
		camera.Begin()
		render.DrawWorld(world, encounter)
		am.Draw(encounter, camera, world)
		camera.End()
		ui.DrawString(fmt.Sprintf("Round %d - Turn: %d", encounter.Round+1, encounter.Turn+1), ui.Rectangle{X: 800, Y: 10, Width: 100, Height: 20}, ui.Black, 20, ui.AlignMiddle)
		render.DrawActions(encounter.ActiveActor(), am.SetActive, am.Active, endTurn)
		window.EndFrame()

		ui.ProcessInput()
		camera.Update()
	}
}

func main() {
	client(nil)
}
