package main

import (
	"anvil/internal/ai"
	"anvil/internal/creature"
	"anvil/internal/encounter"
	"anvil/internal/team"
	"fmt"
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

func main() {
	wizard := creature.New("Wizard", team.Player, 22, []creature.Action{AttackAction{}})
	elf := creature.New("Elf", team.Player, 22, []creature.Action{AttackAction{}})
	orc := creature.New("Orc", team.Enemy, 22, []creature.Action{AttackAction{}})
	goblin := creature.New("Goblin", team.Enemy, 22, []creature.Action{AttackAction{}})
	start := time.Now()
	resultCh := make(chan team.Team)
	go encounter.Play([]*creature.Creature{wizard, elf, orc, goblin}, ai.Simple, resultCh)
	winner := <-resultCh
	fmt.Println("Winner:", winner)
	fmt.Printf("%v elapsed\n", time.Since(start))
}
