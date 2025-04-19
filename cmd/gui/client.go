package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"anvil/cmd/gui/render"
	"anvil/cmd/gui/ui"
	"anvil/internal/demo"
	"anvil/internal/eventbus"
	"anvil/internal/prettyprint"
)

func client(_ net.Conn) {
	log := ui.ScrollText{
		Rect:       ui.Rectangle{X: 600, Y: 40, Width: 650, Height: 580},
		LineHeight: 18 + 4,
		Padding:    4,
		BgColor:    ui.LightGray,
		TextColor:  ui.Black,
		FontSize:   18,
	}

	hub := eventbus.Hub{}
	hub.Subscribe(func(msg eventbus.Message) {
		prettyprint.Print(&log, msg)
	})

	window := render.Window{}
	window.Open()
	defer window.Close()
	ui.Init()
	defer ui.Close()
	world, encounter := demo.Create(&hub)
	camera := render.Camera{}
	camera.Reset(window.Width, window.Height)
	camera.SetPosition(-20, -20)
	am := render.ActionManager{Encounter: encounter, World: world}

	endTurn := func() {
		encounter.EndTurn()
		am.SetActive(nil)
	}

	encounter.Start()
	for !window.ShouldClose() {
		window.StartFrame()
		camera.Begin()
		render.DrawWorld(world, encounter)
		am.Draw(camera)
		camera.End()
		best := encounter.ActiveActor().BestScoredAction()
		if best == nil {
			ui.DrawString("Best Action: End Turn", ui.Rectangle{X: 800, Y: 10, Width: 400, Height: 20}, ui.White, 20, ui.AlignRight)
		} else {
			ui.DrawString(fmt.Sprintf("Best Action: %s", best.Action.Name()), ui.Rectangle{X: 800, Y: 10, Width: 400, Height: 20}, ui.White, 20, ui.AlignRight)
		}
		ui.DrawString(fmt.Sprintf("Round %d - Turn: %d", encounter.Round+1, encounter.Turn+1), ui.Rectangle{X: 800, Y: 10, Width: 00, Height: 20}, ui.White, 20, ui.AlignLeft)
		render.DrawActions(encounter.ActiveActor(), am.SetActive, am.Active, endTurn)
		log.Draw()
		window.EndFrame()
		consumed := ui.ProcessInput()
		if !consumed {
			am.ProcessInput(camera)
		}
		camera.Update()
		ui.Update()
	}
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		// Do cleanup here if needed
		os.Exit(0)
	}()
	client(nil)
}
