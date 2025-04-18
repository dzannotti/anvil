package main

import (
	"net"

	"anvil/cmd/gui/render"
	"anvil/cmd/gui/ui"
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/ruleset"
	"anvil/internal/ruleset/fighter"
	"anvil/internal/ruleset/item/armor"
	"anvil/internal/ruleset/item/weapon"
	"anvil/internal/ruleset/monster/undead/zombie"
	"anvil/internal/tag"
)

func createWorld() *core.World {
	world := core.NewWorld(10, 10)
	walls := make([]grid.Position, 0, 256)
	for x := 0; x < world.Width(); x++ {
		walls = append(walls, grid.Position{X: x, Y: 0}, grid.Position{X: x, Y: world.Height() - 1}, grid.Position{X: 0, Y: x}, grid.Position{X: world.Width() - 1, Y: x})
		if x > 0 && x < world.Width()-2 {
			walls = append(walls, grid.Position{X: world.Width() - x, Y: x})
		}
	}
	for _, p := range walls {
		cell, _ := world.At(p)
		cell.Tile = core.Wall
	}
	return world
}

func createEncounter(hub *eventbus.Hub, world *core.World) *core.Encounter {
	cres := core.Resources{Max: map[tag.Tag]int{
		tags.WalkSpeed: 5,
	}}
	cedric := ruleset.NewPCActor(hub, world, grid.Position{X: 3, Y: 2}, "Cedric", 12, stats.Attributes{Strength: 16, Dexterity: 13, Constitution: 14, Intelligence: 8, Wisdom: 14, Charisma: 10}, stats.Proficiencies{Bonus: 2}, cres)
	cedric.Equip(weapon.NewGreatAxe())
	cedric.Equip(armor.NewChainMail())
	cedric.AddEffect(fighter.NewFightingStyleDefense())
	cedric.AddProficiency(tags.MartialWeapon)
	mob1 := zombie.New(hub, world, grid.Position{X: 7, Y: 6}, "Zombie 1")
	encounter := &core.Encounter{
		Log:    hub,
		World:  world,
		Actors: []*core.Actor{cedric, mob1},
	}
	return encounter
}

func client(_ net.Conn) {
	hub := eventbus.Hub{}
	hub.Subscribe(func(msg eventbus.Message) {})

	window := render.Window{}
	window.Open()
	defer window.Close()
	ui.Init()
	defer ui.Close()
	world := createWorld()
	encounter := createEncounter(&hub, world)
	encounter.Start()
	camera := render.Camera{}
	camera.Reset(window.Width, window.Height)
	camera.SetPosition(-20, -20)
	var activeAction core.Action
	selectAction := func(action core.Action) {
		activeAction = action
	}
	endTurn := func() {
		encounter.EndTurn()
		activeAction = nil
	}
	for !window.ShouldClose() {
		window.StartFrame()
		camera.Begin()
		render.DrawWorld(world)
		camera.End()
		window.EndFrame()
		ui.ProcessInput()
		render.DrawActions(encounter.ActiveActor(), selectAction, activeAction, endTurn)
		camera.Update()
	}
}

func main() {
	client(nil)
}
