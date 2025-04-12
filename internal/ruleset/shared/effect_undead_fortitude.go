package shared

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
)

func NewUndeadFortitudeEffect() *core.Effect {
	fx := &core.Effect{Name: "Undead Fortitude", Priority: core.PriorityLate}

	fx.WithAfterTakeDamage(func(_ *core.Effect, s *core.AfterTakeDamageState) {
		wouldDie := s.Source.HitPoints == 0
		radiant := s.Result.IsDamageType(tags.Radiant)
		if !wouldDie || radiant || s.Result.IsCriticalSuccess() {
			return
		}
		s.Source.Log.Start(core.EffectType, core.EffectEvent{Source: s.Source, Effect: fx})
		defer s.Source.Log.End()
		dc := 5 + s.Result.Value
		st := s.Source.SaveThrow(tags.Constitution, dc)
		if st.Success {
			s.Source.ModifyAttribute(tags.HitPoints, 1, "Undead Fortitude")
		}
	})

	return fx
}
