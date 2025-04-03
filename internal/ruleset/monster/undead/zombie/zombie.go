package zombie

import (
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/ruleset"
)

func New(hub *eventbus.Hub, world *core.World, pos grid.Position, name string) *core.Actor {
	attributes := stats.Attributes{
		Strength:     13,
		Dexterity:    6,
		Constitution: 16,
		Intelligence: 3,
		Wisdom:       6,
		Charisma:     5,
	}
	proficiencies := stats.Proficiencies{Bonus: 2}
	npc := ruleset.NewNPCActor(hub, world, pos, name, 22, attributes, proficiencies)
	return npc
}
