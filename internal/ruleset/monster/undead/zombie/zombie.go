package zombie

import (
	"anvil/internal/core"
	"anvil/internal/core/stats"
	"anvil/internal/core/tags"
	"anvil/internal/eventbus"
	"anvil/internal/expression"
	"anvil/internal/grid"
	"anvil/internal/ruleset/actor"
	"anvil/internal/ruleset/base"
	"anvil/internal/ruleset/shared"
	"anvil/internal/tag"
)

type SlamAttack struct {
	name string
	tags tag.Container
}

func (s SlamAttack) Name() string {
	return s.name
}

func (s SlamAttack) Damage() *expression.Expression {
	expr := expression.FromDamageDice(1, 6, "Slam", s.tags)
	return &expr
}

func (s SlamAttack) Tags() *tag.Container {
	return &s.tags
}

func NewSlamAction(owner *core.Actor) core.Action {
	slam := SlamAttack{
		name: "Slam",
		tags: tag.NewContainer(tags.Bludgeoning),
	}
	return base.NewAttackAction(owner, "Slam", slam, 1, tag.NewContainer(tags.Melee, tags.NaturalWeapon))
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
	npc := actor.NewNPCActor(dispatcher, world, pos, name, 22, attributes, proficiencies, resources)
	npc.AddAction(NewSlamAction(npc))
	npc.AddEffect(shared.NewUndeadFortitudeEffect())
	return npc
}
