package main

import (
	"fmt"
)

type Creature struct {
	Name      string
	HitPoints int
}

func (c *Creature) TakeDamage(damage int) {
	c.HitPoints -= damage
	fmt.Println(c.Name, "took", damage, "damage", c.HitPoints, "remaining")
}

func (c *Creature) IsDead() bool {
	return c.HitPoints == 0
}

type Faction struct {
	Name    string
	Members []*Creature
}

func (f *Faction) Add(creature *Creature) {
	f.Members = append(f.Members, creature)
}

func (f *Faction) IsDead() bool {
	for _, creature := range f.Members {
		if !creature.IsDead() {
			return false
		}
	}
	return true
}

func (f *Faction) Contains(creature *Creature) bool {
	for _, c := range f.Members {
		if c.Name == creature.Name {
			return true
		}
	}
	return false
}

func IsOver(factions []Faction) bool {
	remainingFactions := 0
	for i := range factions {
		if !factions[i].IsDead() {
			remainingFactions++
		}
	}
	return remainingFactions <= 1
}

func AllCreatures(factions []Faction) []*Creature {
	var allCreatures = []*Creature{}
	for _, faction := range factions {
		allCreatures = append(allCreatures, faction.Members...)
	}
	return allCreatures
}

func Winner(factions []Faction) string {
	for _, faction := range factions {
		if !faction.IsDead() {
			return faction.Name
		}
	}
	return "all alive?"
}

func FindEnemies(creature *Creature, factions []Faction) []*Creature {
	var enemies = []*Creature{}
	for i := range factions {
		if factions[i].Contains(creature) {
			continue
		}
		enemies = append(enemies, factions[i].Members...)
	}
	return enemies
}

func FindTarget(enemies []*Creature) *Creature {
	for j := range enemies {
		if !enemies[j].IsDead() {
			return enemies[j]
		}
	}
	return nil
}

func Encounter(factions []Faction) {
	var allCreatures = AllCreatures(factions)
	for !IsOver(factions) {
		for i := range allCreatures {
			var activeCreature = allCreatures[i]
			fmt.Println(activeCreature.Name, "turn")
			if activeCreature.IsDead() {
				fmt.Println(activeCreature.Name, "Skipping dead")
				continue
			}
			enemies := FindEnemies(activeCreature, factions)
			target := FindTarget(enemies)
			if target == nil {
				fmt.Println(activeCreature.Name, "Skipping no target")
				continue
			}
			fmt.Println(activeCreature.Name, "attacks", target.Name)
			target.TakeDamage(5)
			if target.IsDead() {
				fmt.Println(target.Name, "is dead")
			}
			if IsOver(factions) {
				break
			}
		}
	}
	fmt.Println("Winner", Winner(factions))
}

func main() {
	f1 := Faction{
		Name: "Players",
		Members: []*Creature{
			{
				Name:      "Wizard",
				HitPoints: 10,
			},
			{
				Name:      "Elf",
				HitPoints: 10,
			},
		},
	}
	f2 := Faction{
		Name: "Enemies",
		Members: []*Creature{
			{
				Name:      "Orc",
				HitPoints: 10,
			},
			{
				Name:      "Goblin",
				HitPoints: 10,
			},
		},
	}
	Encounter([]Faction{f1, f2})
}
