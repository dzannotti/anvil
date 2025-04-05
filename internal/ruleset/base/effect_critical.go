package base

import "anvil/internal/core"

func NewCritEffect() *core.Effect {
	fx := &core.Effect{Name: "Crit", Priority: core.PriorityLast}

	fx.WithBeforeDamageRoll(func(_ *core.Effect, s *core.BeforeDamageRollState) {
		if *s.Critical {
			s.Expression.DoubleDice("Critical")
		}
	})

	return fx
}
