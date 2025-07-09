package armor

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/tag"

	"github.com/google/uuid"
)

type ChainMail struct {
	archetype string
	id        string
	tags      tag.Container
}

func newChainMailEffect() *core.Effect {
	// TODO: implement item requirements for proficiencies (and maluses)
	fx := &core.Effect{
		Archetype: "chain-mail",
		ID:        uuid.New().String(),
		Name:      "ChainMail",
		Priority:  core.PriorityBaseOverride,
	}

	fx.On(func(s *core.AttributeCalculation) {
		if !s.Attribute.MatchExact(tags.ArmorClass) {
			return
		}

		s.Expression.ReplaceWith(16, "Chain Mail")
	})

	return fx
}

func NewChainMail() *ChainMail {
	return &ChainMail{
		archetype: "chain-mail",
		id:        uuid.New().String(),
		tags:      tag.NewContainer(tags.MediumArmor),
	}
}

func (c *ChainMail) Archetype() string {
	return c.archetype
}

func (c *ChainMail) ID() string {
	return c.id
}

func (c *ChainMail) Name() string {
	return "Chain Mail"
}

func (c *ChainMail) Tags() *tag.Container {
	return &c.tags
}

func (c *ChainMail) OnEquip(a *core.Actor) {
	a.AddEffect(newChainMailEffect())
}
