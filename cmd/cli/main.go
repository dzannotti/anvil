package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"anvil/internal/ai"
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/prettyprint"
	"anvil/internal/ruleset/base"
)

func creature(hub *eventbus.Hub, world *core.World, pos grid.Position, name string, hitPoints int, attributes stats.Attributes, proficiencies stats.Proficiencies) *core.Creature {
	c := &core.Creature{
		Log:           hub,
		Position:      pos,
		World:         world,
		Name:          name,
		HitPoints:     hitPoints,
		MaxHitPoints:  hitPoints,
		Attributes:    attributes,
		Proficiencies: proficiencies,
	}
	world.AddOccupant(pos, c)
	c.AddAction(base.NewAttackAction(c))
	return c
}

func setupWorld(world *core.World) {
	for x := 0; x < world.Width(); x++ {
		cell, _ := world.Navigation.At(grid.Position{X: x, Y: 0})
		cell.SetWalkable(false)
	}
	for x := 0; x < world.Width(); x++ {
		cell, _ := world.Navigation.At(grid.Position{X: x, Y: world.Height() - 1})
		cell.SetWalkable(false)
	}
	for y := 0; y < world.Height(); y++ {
		cell, _ := world.Navigation.At(grid.Position{X: 0, Y: y})
		cell.SetWalkable(false)
	}
	for y := 0; y < world.Height(); y++ {
		cell, _ := world.Navigation.At(grid.Position{X: world.Width() - 1, Y: y})
		cell.SetWalkable(false)
	}
}

func main() {
	hub := eventbus.Hub{}
	hub.Subscribe(func(msg eventbus.Message) {
		prettyprint.Print(os.Stdout, msg)
	})
	world := core.NewWorld(10, 10)
	setupWorld(world)
	players := &core.Team{Name: "Players"}
	enemies := &core.Team{Name: "Enemies"}
	attributes := stats.Attributes{Strength: 10, Dexterity: 11, Constitution: 12, Intelligence: 13, Wisdom: 14, Charisma: 15}
	wizard := creature(&hub, world, grid.Position{X: 1, Y: 1}, "Wizard", 22, attributes, stats.Proficiencies{Bonus: 2})
	fighter := creature(&hub, world, grid.Position{X: 1, Y: 2}, "Fighter", 22, attributes, stats.Proficiencies{Bonus: 2})
	orc := creature(&hub, world, grid.Position{X: 4, Y: 3}, "Orc", 22, attributes, stats.Proficiencies{Bonus: 2})
	goblin := creature(&hub, world, grid.Position{X: 4, Y: 4}, "Goblin", 22, attributes, stats.Proficiencies{Bonus: 2})
	players.AddMember(wizard)
	players.AddMember(fighter)
	enemies.AddMember(orc)
	enemies.AddMember(goblin)
	encounter := core.NewEncounter(&hub, world, []*core.Team{players, enemies})
	gameAI := map[*core.Creature]ai.AI{
		wizard:  ai.NewSimple(encounter, wizard),
		fighter: ai.NewSimple(encounter, fighter),
		orc:     ai.NewSimple(encounter, orc),
		goblin:  ai.NewSimple(encounter, goblin),
	}
	wg := sync.WaitGroup{}
	start := time.Now()
	wg.Add(1)
	go encounter.Play(func(active *core.Creature, wg *sync.WaitGroup) {
		defer wg.Done()
		gameAI[active].Play()
	}, &wg)
	wg.Wait()
	winner, ok := encounter.Winner()
	if !ok {
		fmt.Println("All dead")
		return
	}
	fmt.Println("Winner:", winner.Name)
	fmt.Printf("%v elapsed\n", time.Since(start))
}
