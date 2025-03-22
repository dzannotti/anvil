package main

import (
	"anvil/internal/core"
	"anvil/internal/team"
	"fmt"
	"sync"
	"time"
)

func BestCombatAction(c *core.Creature, creatures []*core.Creature) (CombatAction, error) {
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
	action core.Action
	target *core.Creature
}

type AttackAction struct{}

func (a AttackAction) Cost() int {
	return 1
}

func (a AttackAction) Perform(source *core.Creature, target *core.Creature, wg *sync.WaitGroup) {
	defer wg.Done()
	source.Consume(a.Cost())
	source.Attack(target)
}

func IsOver(creatures []*core.Creature) bool {
	playersAlive := false
	enemiesAlive := false
	for _, c := range creatures {
		if !c.IsDead() {
			if c.FactionID() == team.Player {
				playersAlive = true
			}
			if c.FactionID() == team.Enemy {
				enemiesAlive = true
			}
		}
	}
	return !playersAlive || !enemiesAlive
}

func winner(creatures []*core.Creature) string {
	for i := range creatures {
		if !creatures[i].IsDead() {
			if creatures[i].FactionID() == team.Player {
				return "Player"
			}
			if creatures[i].FactionID() == team.Enemy {
				return "Enemy"
			}
		}
	}
	return "all alive?"
}

func findEnemies(creature *core.Creature, allCreatures []*core.Creature) []*core.Creature {
	var enemies = make([]*core.Creature, 0)
	for i := range allCreatures {
		if allCreatures[i].FactionID() == creature.FactionID() {
			continue
		}
		enemies = append(enemies, allCreatures[i])
	}
	return enemies
}

func findTarget(enemies []*core.Creature) *core.Creature {
	for j := range enemies {
		if !enemies[j].IsDead() {
			return enemies[j]
		}
	}
	return nil
}

func Act(activeCreature *core.Creature, allCreatures []*core.Creature, actWG *sync.WaitGroup) {
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

func Encounter(allCreatures []*core.Creature) {
	turn := 0
	round := 0
	for !IsOver(allCreatures) {
		fmt.Println("Round", round+1)
		for i := range allCreatures {
			var activeCreature = allCreatures[i]
			fmt.Println("Turn", turn+1, activeCreature.Name(), "turn")
			wg := sync.WaitGroup{}
			wg.Add(1)
			activeCreature.StartTurn()
			go Act(activeCreature, allCreatures, &wg)
			wg.Wait()
			turn = turn + 1
			if IsOver(allCreatures) {
				break
			}
		}
		round = round + 1
		turn = 0
	}
	fmt.Println("Winner", winner(allCreatures))
}

func main() {
	wizard := core.NewCreature("Wizard", team.Player, 22, []core.Action{AttackAction{}})
	elf := core.NewCreature("Elf", team.Player, 22, []core.Action{AttackAction{}})
	orc := core.NewCreature("Orc", team.Enemy, 22, []core.Action{AttackAction{}})
	goblin := core.NewCreature("Goblin", team.Enemy, 22, []core.Action{AttackAction{}})
	start := time.Now()
	Encounter([]*core.Creature{wizard, elf, orc, goblin})
	fmt.Printf("%v elapsed\n", time.Since(start))
}
