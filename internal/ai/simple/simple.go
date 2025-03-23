package simple

import (
	"anvil/internal/core/creature"
	"fmt"
	"sync"
)

func bestCombatAction(c *creature.Creature, creatures []*creature.Creature) (CombatAction, error) {
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

func Simple(activeCreature *creature.Creature, allCreatures []*creature.Creature, actWG *sync.WaitGroup) {
	defer actWG.Done()
	if activeCreature.IsDead() {
		return
	}
	for {
		action, err := bestCombatAction(activeCreature, allCreatures)
		if err != nil {
			break
		}
		wg := sync.WaitGroup{}
		wg.Add(1)
		go action.action.Perform(activeCreature, action.target, &wg)
		wg.Wait()
	}
}
