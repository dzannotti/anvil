package ruleset

import "anvil/internal/loader"

func weaponDefinitions() map[string]loader.WeaponDefinition {
	return map[string]loader.WeaponDefinition{
		"dagger": {
			Archetype: "dagger",
			Name:      "Dagger",
			Damage: []loader.DamageData{
				{Formula: "1d4", Kind: "piercing"},
			},
			Tags:  []string{"light", "finesse", "thrown"},
			Reach: 1,
		},
		"greataxe": {
			Archetype: "greataxe",
			Name:      "Great Axe",
			Damage: []loader.DamageData{
				{Formula: "2d6", Kind: "slashing"},
			},
			Tags:  []string{"heavy", "two-handed"},
			Reach: 1,
		},
		"flamingsword": {
			Archetype: "flamingsword",
			Name:      "Flaming Sword",
			Damage: []loader.DamageData{
				{Formula: "1d8", Kind: "slashing"},
				{Formula: "1d6", Kind: "fire"},
			},
			Tags:  []string{"versatile"},
			Reach: 1,
		},
		"zombie_slam": {
			Archetype: "zombie_slam",
			Name:      "Zombie Slam",
			Damage: []loader.DamageData{
				{Formula: "1d6+1", Kind: "bludgeoning"},
			},
			Tags:  []string{"natural"},
			Reach: 1,
		},
	}
}