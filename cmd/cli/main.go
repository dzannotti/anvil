package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"anvil/internal/core"
	"anvil/internal/core/ai"
	"anvil/internal/core/definition"
	"anvil/internal/core/world"
	"anvil/internal/grid"
	"anvil/internal/log"
	"anvil/internal/prettyprint"
	"anvil/internal/ruleset/base"
)

func printLog(event log.Event) {
	prettyprint.Print(os.Stdout, event)
}

func creature(log *log.EventLog, world *world.World, pos grid.Position, name string, hitPoints int) *core.Creature {
	c := core.NewCreature(log, world, pos, name, hitPoints)
	c.AddAction(base.NewAttackAction(c))
	return c
}

func main() {
	log := log.New()
	log.AddCapturer(printLog)
	world := core.NewWorld(10, 10)
	players := core.NewTeam("Players")
	enemies := core.NewTeam("Enemies")
	wizard := creature(log, world, grid.NewPosition(1, 1), "Wizard", 22)
	fighter := creature(log, world, grid.NewPosition(1, 2), "Fighter", 22)
	orc := creature(log, world, grid.NewPosition(4, 3), "Orc", 22)
	goblin := creature(log, world, grid.NewPosition(4, 4), "Goblin", 22)
	players.AddMember(wizard)
	players.AddMember(fighter)
	enemies.AddMember(orc)
	enemies.AddMember(goblin)
	encounter := core.NewEncounter(log, []definition.Team{players, enemies})
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
