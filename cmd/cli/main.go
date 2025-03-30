package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"anvil/internal/core"
	"anvil/internal/core/ai"
	"anvil/internal/core/definition"
	"anvil/internal/grid"
	"anvil/internal/log"
	"anvil/internal/prettyprint"
	"anvil/internal/ruleset/base"
)

func printLog(event log.Event) {
	prettyprint.Print(os.Stdout, event)
}

func creature(log *log.EventLog, world *core.World, pos grid.Position, name string, hitPoints int, attributes core.Attributes, proficiencies core.Proficiencies) *core.Creature {
	c := core.NewCreature(log, world, pos, name, hitPoints, attributes, proficiencies)
	c.AddAction(base.NewAttackAction(c))
	return c
}

func setupWorld(world *core.World) {
	for x := 0; x < world.Width(); x++ {
		cell, _ := world.Navigation().At(grid.NewPosition(x, 0))
		cell.SetWalkable(false)
	}
	for x := 0; x < world.Width(); x++ {
		cell, _ := world.Navigation().At(grid.NewPosition(x, world.Height()-1))
		cell.SetWalkable(false)
	}
	for y := 0; y < world.Height(); y++ {
		cell, _ := world.Navigation().At(grid.NewPosition(0, y))
		cell.SetWalkable(false)
	}
	for y := 0; y < world.Height(); y++ {
		cell, _ := world.Navigation().At(grid.NewPosition(world.Width()-1, y))
		cell.SetWalkable(false)
	}
}

func main() {
	log := log.New()
	log.AddCapturer(printLog)
	world := core.NewWorld(10, 10)
	setupWorld(world)
	players := core.NewTeam("Players")
	enemies := core.NewTeam("Enemies")
	attributes := core.NewAttributes(core.AttributeValues{Strength: 10, Dexterity: 11, Constitution: 12, Intelligence: 13, Wisdom: 14, Charisma: 15})
	wizard := creature(log, world, grid.NewPosition(1, 1), "Wizard", 22, attributes, core.NewProficiencies(2))
	fighter := creature(log, world, grid.NewPosition(1, 2), "Fighter", 22, attributes, core.NewProficiencies(2))
	orc := creature(log, world, grid.NewPosition(4, 3), "Orc", 22, attributes, core.NewProficiencies(2))
	goblin := creature(log, world, grid.NewPosition(4, 4), "Goblin", 22, attributes, core.NewProficiencies(2))
	players.AddMember(wizard)
	players.AddMember(fighter)
	enemies.AddMember(orc)
	enemies.AddMember(goblin)
	encounter := core.NewEncounter(log, world, []definition.Team{players, enemies})
	gameAI := map[definition.Creature]ai.AI{
		wizard:  ai.NewSimple(encounter, wizard),
		fighter: ai.NewSimple(encounter, fighter),
		orc:     ai.NewSimple(encounter, orc),
		goblin:  ai.NewSimple(encounter, goblin),
	}
	wg := sync.WaitGroup{}
	start := time.Now()
	wg.Add(1)
	go encounter.Play(func(active definition.Creature, wg *sync.WaitGroup) {
		defer wg.Done()
		gameAI[active].Play()
	}, &wg)
	wg.Wait()
	fmt.Println("Winner:", encounter.Winner().Name())
	fmt.Printf("%v elapsed\n", time.Since(start))
}
