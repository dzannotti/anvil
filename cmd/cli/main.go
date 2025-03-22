package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Team int

const (
	Player Team = iota + 1
	Enemy
)

func RollDice(sides int) int {
	return rand.Intn(sides) + 1
}

type Creature struct {
	name         string
	factionId    Team
	hitPoints    int
	actionPoints int
	actions      []Action
}

func NewCreature(name string, factionId Team, hitPoints int) *Creature {
	return &Creature{name: name, factionId: factionId, hitPoints: hitPoints, actionPoints: 0, actions: []Action{AttackAction{}}}
}

func (c *Creature) takeDamage(damage int) {
	c.hitPoints = max(0, c.hitPoints-damage)
	fmt.Println(c.name, "took", damage, "damage", c.hitPoints, "remaining")
}

func (c *Creature) Attack(target *Creature) {
	fmt.Println(c.name, "attacks", target.name)
	damage := RollDice(20)
	target.takeDamage(damage)
	if target.isDead() {
		fmt.Println(target.name, "is dead")
	}
}

func (c *Creature) isDead() bool {
	return c.hitPoints == 0
}

func (c *Creature) StartTurn() {
	c.actionPoints = 2
}

func (c *Creature) Consume(cost int) {
	c.actionPoints = max(0, c.actionPoints-cost)
}

func (c *Creature) BestCombatAction(creatures []*Creature) (CombatAction, error) {
	enemies := findEnemies(c, creatures)
	target := findTarget(enemies)
	if target == nil {
		return CombatAction{}, fmt.Errorf("no target")
	}
	action := c.actions[0]
	if c.actionPoints < action.Cost() {
		return CombatAction{}, fmt.Errorf("not enough action points")
	}
	return CombatAction{action: action, target: target}, nil
}

type CombatAction struct {
	action Action
	target *Creature
}

type Action interface {
	Perform(source *Creature, target *Creature, wg *sync.WaitGroup)
	Cost() int
}

type AttackAction struct{}

func (a AttackAction) Cost() int {
	return 1
}

func (a AttackAction) Perform(source *Creature, target *Creature, wg *sync.WaitGroup) {
	defer wg.Done()
	source.Consume(a.Cost())
	source.Attack(target)
}

func IsOver(creatures []*Creature) bool {
	playersAlive := false
	enemiesAlive := false
	for _, c := range creatures {
		if !c.isDead() {
			if c.factionId == Player {
				playersAlive = true
			}
			if c.factionId == Enemy {
				enemiesAlive = true
			}
		}
	}
	return !playersAlive || !enemiesAlive
}

func winner(creatures []*Creature) string {
	for i := range creatures {
		if !creatures[i].isDead() {
			if creatures[i].factionId == Player {
				return "Player"
			}
			if creatures[i].factionId == Enemy {
				return "Enemy"
			}
		}
	}
	return "all alive?"
}

func findEnemies(creature *Creature, allCreatures []*Creature) []*Creature {
	var enemies = make([]*Creature, 0)
	for i := range allCreatures {
		if allCreatures[i].factionId == creature.factionId {
			continue
		}
		enemies = append(enemies, allCreatures[i])
	}
	return enemies
}

func findTarget(enemies []*Creature) *Creature {
	for j := range enemies {
		if !enemies[j].isDead() {
			return enemies[j]
		}
	}
	return nil
}

func Act(activeCreature *Creature, allCreatures []*Creature, actWG *sync.WaitGroup) {
	defer actWG.Done()
	if activeCreature.isDead() {
		fmt.Println(activeCreature.name, "cannot act because dead")
		return
	}
	for {
		action, err := activeCreature.BestCombatAction(allCreatures)
		if err != nil {
			fmt.Println(activeCreature.name, "cannot act: no action")
			break
		}
		wg := sync.WaitGroup{}
		wg.Add(1)
		go action.action.Perform(activeCreature, action.target, &wg)
		wg.Wait()
	}
}

func Encounter(allCreatures []*Creature) {
	turn := 0
	round := 0
	for !IsOver(allCreatures) {
		fmt.Println("Round", round+1)
		for i := range allCreatures {
			var activeCreature = allCreatures[i]
			fmt.Println("Turn", turn+1, activeCreature.name, "turn")
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
	wizard := NewCreature("Wizard", Player, 22)
	elf := NewCreature("Elf", Player, 22)
	orc := NewCreature("Orc", Enemy, 22)
	goblin := NewCreature("Goblin", Enemy, 22)
	start := time.Now()
	Encounter([]*Creature{wizard, elf, orc, goblin})
	fmt.Printf("%v elapsed\n", time.Since(start))
}
