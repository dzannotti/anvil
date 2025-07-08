package base

import "anvil/internal/core"

func NewCritEffect() *core.Effect {
	fx := &core.Effect{Name: "Crit", Priority: core.PriorityLate}

	fx.On(func(s *core.BeforeDamageRollState) {
		if s.Expression.IsCriticalSuccess() {
			s.Expression.DoubleDice("Critical")
		}
	})

	return fx
}
