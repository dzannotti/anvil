package basic

import "anvil/internal/core"

func NewDeathEffect() *core.Effect {
	fx := &core.Effect{Name: "Death", Priority: core.PriorityLast}

	fx.On(func(s *core.PostTakeDamage) {
		if s.Source.HitPoints == 0 {
			s.Source.Die()
		}
	})
	return fx
}
