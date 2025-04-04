package armor

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
)

type ChainMail struct{}

func newChainMailEffect() *core.Effect {
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
	return &ChainMail{}
}

func (c *ChainMail) Name() string {
	return "Chain Mail"
}

func (c *ChainMail) OnEquip(a *core.Actor) {
	a.AddEffect(newChainMailEffect())
}
