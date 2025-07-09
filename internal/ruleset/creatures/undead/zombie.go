package undead

import (
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/ruleset/factories"
	"anvil/internal/tag"
)

// RegistryReader provides read-only access to the registry
type RegistryReader interface {
	NewAction(archetype string, owner *core.Actor, options map[string]interface{}) (core.Action, error)
	NewEffect(archetype string, options map[string]interface{}) (*core.Effect, error)
	NewItem(archetype string, options map[string]interface{}) (core.Item, error)
	NewCreature(archetype string, options map[string]interface{}) (*core.Actor, error)
	HasAction(archetype string) bool
	HasEffect(archetype string) bool
	HasItem(archetype string) bool
	HasCreature(archetype string) bool
}

func New(registry RegistryReader, dispatcher *eventbus.Dispatcher, world *core.World, pos grid.Position, name string) *core.Actor {
	attributes := stats.Attributes{
		Strength:     13,
		Dexterity:    6,
		Constitution: 16,
		Intelligence: 3,
		Wisdom:       6,
		Charisma:     5,
	}
	proficiencies := stats.Proficiencies{Bonus: 2}
	resources := core.Resources{Max: map[tag.Tag]int{
		tags.WalkSpeed: 4,
	}}
	npc := factories.NewNPCActor(registry, dispatcher, world, pos, name, 22, attributes, proficiencies, resources)

	slamAction, err := registry.NewAction("slam", npc, nil)
	if err != nil {
		panic("failed to create slam action: " + err.Error())
	}
	npc.AddAction(slamAction)

	undeadFortitude, err := registry.NewEffect("undead-fortitude", nil)
	if err != nil {
		panic("failed to create undead fortitude effect: " + err.Error())
	}
	npc.AddEffect(undeadFortitude)

	return npc
}
