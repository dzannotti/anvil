package seed

import "anvil/internal/loader"

func ZombieDefinition(name string) loader.ActorDefinition {
	return loader.ActorDefinition{
		Name:         name,
		Team:         "enemies",
		HitPoints:    22,
		MaxHitPoints: 22,
		Attributes: loader.AttributesDefinition{
			Strength:     13,
			Dexterity:    6,
			Constitution: 16,
			Intelligence: 3,
			Wisdom:       6,
			Charisma:     5,
		},
		Proficiencies: loader.ProficienciesDefinition{
			Skills: []string{},
			Bonus:  2,
		},
		Resources: loader.ResourcesDefinition{
			WalkSpeed: 4,
		},
	}
}

func ZombieSlamDefinition() loader.MeleeActionDefinition {
	return loader.MeleeActionDefinition{
		Name:          "Zombie Slam",
		Cost:          map[string]int{"action": 1},
		Tags:          []string{"attack", "natural"},
		Reach:         1,
		DamageFormula: "1d6",
		DamageType:    "bludgeoning",
	}
}
