package basic

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/expression"
)

func NewUndeadFortitudeEffect() *core.Effect {
	fx := &core.Effect{Name: "Undead Fortitude", Priority: core.PriorityLate}

	fx.On(func(s *core.PostTakeDamage) {
		wouldDie := s.Source.HitPoints == 0
		radiant := s.Result.HasDamageType(expression.DamageRadiant)
		if !wouldDie || radiant || s.Result.IsCriticalSuccess() {
			return
		}
		s.Source.Dispatcher.Begin(core.EffectEvent{Source: s.Source, Effect: fx})
		defer s.Source.Dispatcher.End()
		dc := 5 + s.Result.Value
		st := s.Source.SaveThrow(tags.AttributeConstitution, dc)
		if st.Success {
			s.Source.ModifyAttribute(tags.ActorHitPoints, 1, "Undead Fortitude")
		}
	})

	return fx
}
