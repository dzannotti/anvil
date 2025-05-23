package armor

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/tag"
)

type ChainMail struct{ tags tag.Container }

func newChainMailEffect() *core.Effect {
	// TODO: implement item requirements for proficiencies (and maluses)
	fx := &core.Effect{Name: "ChainMail", Priority: core.PriorityBaseOverride}

	fx.WithAttributeCalculation(func(_ *core.Effect, s *core.AttributeCalculationState) {
		if !s.Attribute.MatchExact(tags.ArmorClass) {
			return
		}
		s.Expression.ReplaceWith(16, "Chain Mail")
	})

	return fx
}

func NewChainMail() *ChainMail {
	return &ChainMail{tags: tag.ContainerFromTag(tags.MediumArmor)}
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
