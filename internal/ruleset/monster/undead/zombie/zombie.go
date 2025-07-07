package zombie

import (
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/grid"
	"anvil/internal/ruleset/actor"
	"anvil/internal/ruleset/base"
	"anvil/internal/ruleset/shared"
	"anvil/internal/tag"
)

func NewSlamAction(owner *core.Actor) core.Action {
	return base.NewAttackAction(owner, "Slam", []core.DamageSource{
		{Times: 1, Sides: 6, Source: "Slam", Tags: tag.NewContainer(tags.Bludgeoning)},
	}, 1, tag.NewContainer(tags.Melee, tags.NaturalWeapon))
}

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
	resources := core.Resources{Max: map[tag.Tag]int{
		tags.WalkSpeed: 4,
	}}
	npc := actor.NewNPCActor(hub, world, pos, name, 22, attributes, proficiencies, resources)
	npc.AddAction(NewSlamAction(npc))
	npc.AddEffect(shared.NewUndeadFortitudeEffect())
	return npc
}
