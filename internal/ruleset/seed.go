package ruleset

import (
	"anvil/internal/core"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/loader"
	"anvil/internal/ruleset/basic"
	"anvil/internal/ruleset/seed"
)

func SeedRegistry(registry *Registry) {
	registerBasicActions(registry)
	registerBasicEffects(registry)
	registerSharedEffects(registry)
	registerClassEffects(registry)
	registerItems(registry)
	registerActors(registry)
}

func registerBasicActions(registry *Registry) {
	registry.RegisterAction("move", func(owner *core.Actor, _ map[string]interface{}) core.Action {
		return basic.NewMoveAction(owner)
	})

	registry.RegisterAction("melee", func(owner *core.Actor, options map[string]interface{}) core.Action {
		def, ok := options["definition"].(loader.MeleeActionDefinition)
		if !ok {
			panic("melee action requires MeleeActionDefinition")
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

	weaponDefs := seed.WeaponDefinitions()
	for archetype, weaponDef := range weaponDefs {
		def := weaponDef
		registry.RegisterItem(archetype, func(_ map[string]interface{}) core.Item {
			return basic.NewWeaponFromDefinition(def)
		})
	}
}

func newZombie(registry *Registry, dispatcher *eventbus.Dispatcher, world *core.World, pos grid.Position, name string) *core.Actor {
	definition := seed.ZombieDefinition(name)
	npc := registry.CreateActorFromDefinition(dispatcher, world, pos, definition)

	// Create zombie slam action from definition
	slamDef := seed.ZombieSlamDefinition()
	slamAction := registry.NewAction("melee", npc, map[string]interface{}{
		"definition": slamDef,
	})
	npc.AddAction(slamAction)

	npc.AddEffect(registry.NewEffect("undead-fortitude", nil))
	return npc
}

func registerActors(registry *Registry) {
	registry.RegisterActor("zombie", func(options map[string]interface{}) *core.Actor {
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
		actions: make(map[string]ActionFactory),
		effects: make(map[string]EffectFactory),
		items:   make(map[string]ItemFactory),
		actors:  make(map[string]ActorFactory),
	}
	SeedRegistry(registry)
	return registry
}
