package basic

import "anvil/internal/core"

func NewCritEffect() *core.Effect {
	fx := &core.Effect{Name: "Crit", Priority: core.PriorityLate}

	fx.On(func(s *core.PreDamageRoll) {
		if s.Critical {
			s.Expression.DoubleDice("Critical")
		}
	})

	return fx
}
