package undead

import (
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/expression"
	"anvil/internal/grid"
	"anvil/internal/ruleset/actions/basic"
	effectsShared "anvil/internal/ruleset/effects/shared"
	"anvil/internal/ruleset/factories"
	"anvil/internal/tag"
)

func NewSlamAction(owner *core.Actor) core.Action {
	damage := expression.FromDamageDice(1, 6, "Slam", tag.NewContainer(tags.Bludgeoning))
	slam := basic.NewNaturalWeapon("Slam", "slam", damage, tag.NewContainer(tags.Bludgeoning))
	cost := map[tag.Tag]int{tags.Action: 1}
	return basic.NewMeleeAction(owner, "Slam", slam, 1, tag.NewContainer(tags.Melee, tags.NaturalWeapon), cost)
}

func New(dispatcher *eventbus.Dispatcher, world *core.World, pos grid.Position, name string) *core.Actor {
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
	npc := factories.NewNPCActor(dispatcher, world, pos, name, 22, attributes, proficiencies, resources)
	npc.AddAction(NewSlamAction(npc))
	npc.AddEffect(effectsShared.NewUndeadFortitudeEffect())
	return npc
}
