package ruleset

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/expression"
	"anvil/internal/grid"
	actionsBasic "anvil/internal/ruleset/actions/basic"
	actionsShared "anvil/internal/ruleset/actions/shared"
	creaturesUndead "anvil/internal/ruleset/creatures/undead"
	effectsBasic "anvil/internal/ruleset/effects/basic"
	effectsFighter "anvil/internal/ruleset/effects/classes/fighter"
	effectsShared "anvil/internal/ruleset/effects/shared"
	itemsArmor "anvil/internal/ruleset/items/armor"
	itemsWeapons "anvil/internal/ruleset/items/weapons"
	"anvil/internal/tag"
)

// SeedRegistry populates the registry with all available archetypes
// Registration order is important: dependencies must be registered first
func SeedRegistry(registry *Registry) {
	// Register basic actions first (no dependencies)
	registerBasicActions(registry)

	// Register shared actions
	registerSharedActions(registry)

	// Register basic effects
	registerBasicEffects(registry)

	// Register shared effects
	registerSharedEffects(registry)

	// Register class-specific effects
	registerClassEffects(registry)

	// Register items (may depend on actions/effects)
	registerItems(registry)

	// Register creatures last (depends on actions/effects/items)
	registerCreatures(registry)
}

// registerBasicActions registers fundamental D&D actions
func registerBasicActions(registry *Registry) {
	registry.RegisterAction("move", func(owner *core.Actor, _ map[string]interface{}) core.Action {
		return actionsBasic.NewMoveAction(owner)
	})

	// Natural weapon actions
	registry.RegisterAction("slam", func(owner *core.Actor, _ map[string]interface{}) core.Action {
		damage := expression.FromDamageDice(1, 6, "Slam", tag.NewContainer(tags.Bludgeoning))
		slam := actionsBasic.NewNaturalWeapon("Slam", "slam", damage, tag.NewContainer(tags.Bludgeoning))
		cost := map[tag.Tag]int{tags.Action: 1}
		return actionsBasic.NewMeleeAction(owner, "Slam", slam, 1, tag.NewContainer(tags.Melee, tags.NaturalWeapon), cost)
	})

	// TODO: Add other basic actions as they're refactored
	// registry.RegisterAction("melee-attack", func(owner *core.Actor, options map[string]interface{}) core.Action {
	//     return basic.NewMeleeAction(owner, ...)
	// })
}

// registerSharedActions registers reusable actions
func registerSharedActions(registry *Registry) {
	registry.RegisterAction("fireball", func(owner *core.Actor, _ map[string]interface{}) core.Action {
		return actionsShared.NewFireballAction(owner)
	})
}

// registerBasicEffects registers fundamental D&D effects
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

// registerSharedEffects registers reusable effects
func registerSharedEffects(registry *Registry) {
	registry.RegisterEffect("undead-fortitude", func(_ map[string]interface{}) *core.Effect {
		return effectsShared.NewUndeadFortitudeEffect()
	})
}

// registerClassEffects registers class-specific effects
func registerClassEffects(registry *Registry) {
	// Fighter effects
	registry.RegisterEffect("fighting-style-defense", func(_ map[string]interface{}) *core.Effect {
		return effectsFighter.NewFightingStyleDefense()
	})
}

// registerItems registers items (weapons, armor, etc.)
func registerItems(registry *Registry) {
	// Armor
	registry.RegisterItem("chainmail", func(_ map[string]interface{}) core.Item {
		return itemsArmor.NewChainMail()
	})

	// Weapons
	registry.RegisterItem("dagger", func(_ map[string]interface{}) core.Item {
		return itemsWeapons.NewDagger()
	})

	registry.RegisterItem("greataxe", func(_ map[string]interface{}) core.Item {
		return itemsWeapons.NewGreatAxe()
	})
}

// registerCreatures registers creature archetypes
func registerCreatures(registry *Registry) {
	registry.RegisterCreature("zombie", func(options map[string]interface{}) *core.Actor {
		// Extract required parameters from options
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

// InitializeDefaultRegistry seeds the global registry with all archetypes
func InitializeDefaultRegistry() {
	SeedRegistry(DefaultRegistry)
}
