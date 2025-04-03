package base

import "anvil/internal/core"

func NewDeathEffect(_ *core.Actor) *core.Effect {
	fx := &core.Effect{Name: "Death", Priority: core.PriorityLast}

	fx.WithAfterTakeDamage(func(_ *core.Effect, s *core.AfterTakeDamageState) {
		if s.Source.HitPoints == 0 {
			s.Source.Die()
		}
	})
	return fx
}
