package main

import (
	"fmt"
	"os"
	"time"

	"anvil/internal/ai"
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/prettyprint"
	"anvil/internal/ruleset"
)

func setupWorld(world *core.World) {
	for x := 0; x < world.Width(); x++ {
		cell, _ := world.Navigation.At(grid.Position{X: x, Y: 0})
		cell.Walkable = false
	}
	for x := 0; x < world.Width(); x++ {
		cell, _ := world.Navigation.At(grid.Position{X: x, Y: world.Height() - 1})
		cell.Walkable = false
	}
	for y := 0; y < world.Height(); y++ {
		cell, _ := world.Navigation.At(grid.Position{X: 0, Y: y})
		cell.Walkable = false
	}
	for y := 0; y < world.Height(); y++ {
		cell, _ := world.Navigation.At(grid.Position{X: world.Width() - 1, Y: y})
		cell.Walkable = false
	}
}

func main() {
	hub := eventbus.Hub{}
	hub.Subscribe(func(msg eventbus.Message) {
		prettyprint.Print(os.Stdout, msg)
	})
	world := core.NewWorld(10, 10)
	setupWorld(world)
	attributes := stats.Attributes{Strength: 10, Dexterity: 11, Constitution: 12, Intelligence: 13, Wisdom: 14, Charisma: 15}
	wizard := ruleset.NewPCActor(&hub, world, grid.Position{X: 1, Y: 1}, "Wizard", 22, attributes, stats.Proficiencies{Bonus: 2})
	fighter := ruleset.NewPCActor(&hub, world, grid.Position{X: 1, Y: 2}, "Fighter", 22, attributes, stats.Proficiencies{Bonus: 2})
	orc := ruleset.NewNPCActor(&hub, world, grid.Position{X: 4, Y: 3}, "Orc", 22, attributes, stats.Proficiencies{Bonus: 2})
	goblin := ruleset.NewNPCActor(&hub, world, grid.Position{X: 4, Y: 4}, "Goblin", 22, attributes, stats.Proficiencies{Bonus: 2})
	encounter := &core.Encounter{
		Log:    &hub,
		World:  world,
		Actors: []*core.Actor{wizard, fighter, orc, goblin},
	}
	gameAI := map[*core.Actor]ai.AI{
		wizard:  &ai.Simple{Encounter: encounter, Owner: wizard},
		fighter: &ai.Simple{Encounter: encounter, Owner: fighter},
		orc:     &ai.Simple{Encounter: encounter, Owner: orc},
		goblin:  &ai.Simple{Encounter: encounter, Owner: goblin},
	}
	start := time.Now()
	winner := encounter.Play(func(active *core.Actor) {
		gameAI[active].Play()
	})
	total := time.Since(start)
	if len(winner) == 0 {
		fmt.Println("All dead")
		return
	}
	fmt.Println("Winner:", winner)
	fmt.Printf("%v elapsed\n", total)
}
