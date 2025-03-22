package main

import (
	"anvil/internal/creature"
	"anvil/internal/encounter"
	"anvil/internal/team"
	"fmt"
	"sync"
	"time"
)

func BestCombatAction(c *creature.Creature, creatures []*creature.Creature) (CombatAction, error) {
	enemies := findEnemies(c, creatures)
	target := findTarget(enemies)
	if target == nil {
		return CombatAction{}, fmt.Errorf("no target")
	}
	action := c.Actions()[0]
	if c.ActionPoints() < action.Cost() {
		return CombatAction{}, fmt.Errorf("not enough action points")
	}
	return CombatAction{action: action, target: target}, nil
}

type CombatAction struct {
	action creature.Action
	target *creature.Creature
}

type AttackAction struct{}

func (a AttackAction) Cost() int {
	return 1
}

func (a AttackAction) Perform(source *creature.Creature, target *creature.Creature, wg *sync.WaitGroup) {
	defer wg.Done()
	source.Consume(a.Cost())
	source.Attack(target)
}

func findEnemies(activeCreature *creature.Creature, allCreatures []*creature.Creature) []*creature.Creature {
	var enemies = make([]*creature.Creature, 0)
	for i := range allCreatures {
		if allCreatures[i].FactionID() == activeCreature.FactionID() {
			continue
		}
		enemies = append(enemies, allCreatures[i])
	}
	return enemies
}

func findTarget(enemies []*creature.Creature) *creature.Creature {
	for j := range enemies {
		if !enemies[j].IsDead() {
			return enemies[j]
		}
	}
	return nil
}

func Act(activeCreature *creature.Creature, allCreatures []*creature.Creature, actWG *sync.WaitGroup) {
	defer actWG.Done()
	if activeCreature.IsDead() {
		fmt.Println(activeCreature.Name(), "cannot act because dead")
		return
	}
	for {
		action, err := BestCombatAction(activeCreature, allCreatures)
		if err != nil {
			fmt.Println(activeCreature.Name(), "cannot act: no action")
			break
		}
		wg := sync.WaitGroup{}
		wg.Add(1)
		go action.action.Perform(activeCreature, action.target, &wg)
		wg.Wait()
	}
}

func main() {
	wizard := creature.New("Wizard", team.Player, 22, []creature.Action{AttackAction{}})
	elf := creature.New("Elf", team.Player, 22, []creature.Action{AttackAction{}})
	orc := creature.New("Orc", team.Enemy, 22, []creature.Action{AttackAction{}})
	goblin := creature.New("Goblin", team.Enemy, 22, []creature.Action{AttackAction{}})
	start := time.Now()
	resultCh := make(chan team.Team)
	go encounter.Play([]*creature.Creature{wizard, elf, orc, goblin}, Act, resultCh)
	winner := <-resultCh
	fmt.Println("Winner:", winner)
	fmt.Printf("%v elapsed\n", time.Since(start))
}
