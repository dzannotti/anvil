package ruleset

import (
	"log"
	"path/filepath"
	"runtime"

	"anvil/internal/core"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	actionsBasic "anvil/internal/ruleset/actions/basic"
	actionsShared "anvil/internal/ruleset/actions/shared"
	creaturesUndead "anvil/internal/ruleset/creatures/undead"
	effectsBasic "anvil/internal/ruleset/effects/basic"
	effectsFighter "anvil/internal/ruleset/effects/classes/fighter"
	effectsShared "anvil/internal/ruleset/effects/shared"
	itemsArmor "anvil/internal/ruleset/items/armor"
	"anvil/internal/ruleset/loader"
)

func SeedRegistry(registry *Registry) {
	registerBasicActions(registry)
	registerSharedActions(registry)
	registerBasicEffects(registry)
	registerSharedEffects(registry)
	registerClassEffects(registry)
	registerItems(registry)
	registerCreatures(registry)
}

func registerBasicActions(registry *Registry) {
	registry.RegisterAction("move", func(owner *core.Actor, _ map[string]interface{}) core.Action {
		return actionsBasic.NewMoveAction(owner)
	})
}

func registerSharedActions(registry *Registry) {
	registry.RegisterAction("fireball", func(owner *core.Actor, _ map[string]interface{}) core.Action {
		return actionsShared.NewFireballAction(owner)
	})
}

func registerBasicEffects(registry *Registry) {
	registry.RegisterEffect("critical", func(_ map[string]interface{}) *core.Effect {
		return effectsBasic.NewCritEffect()
	})

	registry.RegisterEffect("death", func(_ map[string]interface{}) *core.Effect {
		return effectsBasic.NewDeathEffect()
	})

	registry.RegisterEffect("death-saving-throw", func(_ map[string]interface{}) *core.Effect {
		return effectsBasic.NewDeathSavingThrowEffect()
	})

	registry.RegisterEffect("attack-of-opportunity", func(_ map[string]interface{}) *core.Effect {
		return effectsBasic.NewAttackOfOpportunityEffect()
	})

	registry.RegisterEffect("proficiency-modifier", func(_ map[string]interface{}) *core.Effect {
		return effectsBasic.NewProficiencyModifierEffect()
	})

	registry.RegisterEffect("attribute-modifier", func(_ map[string]interface{}) *core.Effect {
		return effectsBasic.NewAttributeModifierEffect()
	})
}

func registerSharedEffects(registry *Registry) {
	registry.RegisterEffect("undead-fortitude", func(_ map[string]interface{}) *core.Effect {
		return effectsShared.NewUndeadFortitudeEffect()
	})
}

func registerClassEffects(registry *Registry) {
	registry.RegisterEffect("fighting-style-defense", func(_ map[string]interface{}) *core.Effect {
		return effectsFighter.NewFightingStyleDefense()
	})
}

func registerItems(registry *Registry) {
	registry.RegisterItem("chainmail", func(_ map[string]interface{}) core.Item {
		return itemsArmor.NewChainMail()
	})

	_, file, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(file)))
	dataDir := filepath.Join(projectRoot, "data")

	weaponFactories, err := loader.LoadWeapons(dataDir)
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

		return creaturesUndead.New(registry, dispatcher, world, pos, name)
	})
}

func InitializeDefaultRegistry() {
	SeedRegistry(DefaultRegistry)
}

func NewSeededRegistry() *Registry {
	registry := NewRegistry()
	SeedRegistry(registry)
	return registry
}
