package ruleset

import (
	"log"
	"path/filepath"
	"runtime"

	"anvil/internal/core"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/loader"
	"anvil/internal/ruleset/basic"
	rulesetLoader "anvil/internal/ruleset/loader"
)

func SeedRegistry(registry *Registry) {
	registerBasicActions(registry)
	registerBasicEffects(registry)
	registerSharedEffects(registry)
	registerClassEffects(registry)
	registerItems(registry)
	registerCreatures(registry)
}

func registerBasicActions(registry *Registry) {
	registry.RegisterAction("move", func(owner *core.Actor, _ map[string]interface{}) core.Action {
		return basic.NewMoveAction(owner)
	})

	registry.RegisterAction("melee", func(owner *core.Actor, options map[string]interface{}) core.Action {
		def, ok := options["definition"].(loader.ActionDefinition)
		if !ok {
			panic("melee action requires ActionDefinition")
		}
		return basic.NewMeleeActionFromDefinition(owner, def)
	})
}

func registerBasicEffects(registry *Registry) {
	registry.RegisterEffect("critical", func(_ map[string]interface{}) *core.Effect {
		return basic.NewCritEffect()
	})

	registry.RegisterEffect("death", func(_ map[string]interface{}) *core.Effect {
		return basic.NewDeathEffect()
	})

	registry.RegisterEffect("death-saving-throw", func(_ map[string]interface{}) *core.Effect {
		return basic.NewDeathSavingThrowEffect()
	})

	registry.RegisterEffect("attack-of-opportunity", func(_ map[string]interface{}) *core.Effect {
		return basic.NewAttackOfOpportunityEffect()
	})

	registry.RegisterEffect("proficiency-modifier", func(_ map[string]interface{}) *core.Effect {
		return basic.NewProficiencyModifierEffect()
	})

	registry.RegisterEffect("attribute-modifier", func(_ map[string]interface{}) *core.Effect {
		return basic.NewAttributeModifierEffect()
	})
}

func registerSharedEffects(registry *Registry) {
	registry.RegisterEffect("undead-fortitude", func(_ map[string]interface{}) *core.Effect {
		return basic.NewUndeadFortitudeEffect()
	})
}

func registerClassEffects(registry *Registry) {
	registry.RegisterEffect("fighting-style-defense", func(_ map[string]interface{}) *core.Effect {
		return basic.NewFightingStyleDefense()
	})
}

func registerItems(registry *Registry) {
	registry.RegisterItem("chainmail", func(_ map[string]interface{}) core.Item {
		return basic.NewChainMail()
	})

	_, file, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(file)))
	dataDir := filepath.Join(projectRoot, "data")

	weaponFactories, err := rulesetLoader.LoadWeapons(dataDir)
	if err != nil {
		log.Fatalf("Failed to load weapons: %v", err)
	}

	for archetype, factory := range weaponFactories {
		f := factory
		registry.RegisterItem(archetype, func(_ map[string]interface{}) core.Item {
			return f()
		})
	}
}

func zombieDefinition(name string) loader.ActorDefinition {
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

func zombieSlamDefinition() loader.ActionDefinition {
	return loader.ActionDefinition{
		Name:      "Zombie Slam",
		Archetype: "melee",
		Cost:      map[string]int{"action": 1},
		Tags:      []string{"attack", "natural"},
		MeleeConfig: &loader.MeleeActionConfig{
			Reach:         1,
			DamageFormula: "1d6",
			DamageType:    "bludgeoning",
		},
	}
}

func newZombie(registry *Registry, dispatcher *eventbus.Dispatcher, world *core.World, pos grid.Position, name string) *core.Actor {
	definition := zombieDefinition(name)
	npc := registry.CreateActorFromDefinition(dispatcher, world, pos, definition)
	
	// Create zombie slam action from definition
	slamDef := zombieSlamDefinition()
	slamAction := registry.NewAction("melee", npc, map[string]interface{}{
		"definition": slamDef,
	})
	npc.AddAction(slamAction)
	
	npc.AddEffect(registry.NewEffect("undead-fortitude", nil))
	return npc
}

func registerCreatures(registry *Registry) {
	registry.RegisterCreature("zombie", func(options map[string]interface{}) *core.Actor {
		dispatcher, ok := options["dispatcher"].(*eventbus.Dispatcher)
		if !ok {
			panic("zombie creation requires dispatcher")
		}

		world, ok := options["world"].(*core.World)
		if !ok {
			panic("zombie creation requires world")
		}

		pos, ok := options["position"].(grid.Position)
		if !ok {
			panic("zombie creation requires position")
		}

		name, ok := options["name"].(string)
		if !ok {
			name = "Zombie" // Default name
		}

		return newZombie(registry, dispatcher, world, pos, name)
	})
}

func NewRegistry() *Registry {
	registry := &Registry{
		actions:   make(map[string]ActionFactory),
		effects:   make(map[string]EffectFactory),
		items:     make(map[string]ItemFactory),
		creatures: make(map[string]CreatureFactory),
	}
	SeedRegistry(registry)
	return registry
}
