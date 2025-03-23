package main

import (
	"anvil/internal/ai"
	"anvil/internal/core/creature"
	"anvil/internal/core/encounter"
	"anvil/internal/core/team"
	"anvil/internal/log"
	"anvil/internal/prettyprint"
	"fmt"
	"os"
	"sync"
	"time"
)

type AttackAction struct{}

func (a AttackAction) Cost() int {
	return 1
}

func (a AttackAction) Perform(source *creature.Creature, target *creature.Creature, wg *sync.WaitGroup) {
	defer wg.Done()
	source.Consume(a.Cost())
	source.Attack(target)
}

func printLog(event log.Event) {
	prettyprint.Print(os.Stdout, event)
}

func main() {
	log := log.New()
	log.AddCapturer(printLog)
	wizard := creature.New(log, "Wizard", team.Player, 22, []creature.Action{AttackAction{}})
	elf := creature.New(log, "Elf", team.Player, 22, []creature.Action{AttackAction{}})
	orc := creature.New(log, "Orc", team.Enemy, 22, []creature.Action{AttackAction{}})
	goblin := creature.New(log, "Goblin", team.Enemy, 22, []creature.Action{AttackAction{}})
	start := time.Now()
	resultCh := make(chan team.Team)
	go encounter.Play(log, []*creature.Creature{wizard, elf, orc, goblin}, ai.Simple, resultCh)
	winner := <-resultCh
	fmt.Println("Winner:", winner)
	fmt.Printf("%v elapsed\n", time.Since(start))
}
