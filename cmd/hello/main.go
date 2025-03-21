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
	Members []Creature
}

func (f *Faction) Add(creature Creature) {
	f.Members = append(f.Members, creature)
}

func (f *Faction) IsDead() bool {
	for _, creature := range f.Members {
		if creature.IsDead() {
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

func Winner(factions []Faction) string {
	for _, faction := range factions {
		if !faction.IsDead() {
			return faction.Name
		}
	}
	return ""
}

func Encounter(factions []Faction) {
	var allCreatures = []Creature{}
	for _, faction := range factions {
		allCreatures = append(allCreatures, faction.Members...)
	}
	for !IsOver(factions) {
		for i := range allCreatures {
			fmt.Println(allCreatures[i].Name, "turn")
			if allCreatures[i].IsDead() {
				continue
			}
			allCreatures[(i+1)%len(allCreatures)].TakeDamage(5)
		}
	}
	fmt.Println("Winner", Winner(factions))
}

func main() {
	f1 := Faction{
		Name: "Players",
		Members: []Creature{
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
		Members: []Creature{
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
