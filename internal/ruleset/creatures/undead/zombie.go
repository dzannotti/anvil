package undead

import (
	"anvil/internal/core"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/loader"
)

type RegistryReader interface {
	NewAction(archetype string, owner *core.Actor, options map[string]interface{}) core.Action
	NewEffect(archetype string, options map[string]interface{}) *core.Effect
	NewItem(archetype string, options map[string]interface{}) core.Item
	NewCreature(archetype string, options map[string]interface{}) *core.Actor
	CreateActorFromDefinition(dispatcher *eventbus.Dispatcher, world *core.World, position grid.Position, definition loader.ActorDefinition) *core.Actor
}

func New(registry RegistryReader, dispatcher *eventbus.Dispatcher, world *core.World, pos grid.Position, name string) *core.Actor {
	definition := ZombieDefinition(name)
	npc := registry.CreateActorFromDefinition(dispatcher, world, pos, definition)
	npc.Equip(registry.NewItem("zombie_slam", nil))
	npc.AddEffect(registry.NewEffect("undead-fortitude", nil))
	return npc
}

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
