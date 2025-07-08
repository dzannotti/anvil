package base

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
)

//nolint:funlen,gocognit,cyclop // reason: this function is long for a valid reason
func NewDeathSavingThrowEffect() *core.Effect {
	fx := &core.Effect{Name: "Death Saving Throw", Priority: core.PriorityLast}
	success := 0
	failures := 0

	reset := func() {
		success = 0
		failures = 0
	}

	stabilize := func(src *core.Actor) {
		reset()
		src.RemoveCondition(tags.Unconscious, fx)
		src.AddCondition(tags.Stable, fx)
	}

	checkStatus := func(src *core.Actor) bool {
		if failures >= 3 {
			src.Die()
			return true
		}
		if success >= 3 {
			stabilize(src)
			return true
		}
		return false
	}

	fx.On(func(s *core.AttributeChangedState) {
		if !s.Attribute.MatchExact(tags.HitPoints) {
			return
		}
		if s.Value != 0 {
			return
		}
		reset()
		s.Source.RemoveCondition(tags.Stable, nil)
		s.Source.RemoveCondition(tags.Unconscious, nil)
	})

	fx.On(func(s *core.ConditionChangedState) {
		if !s.Condition.Match(tags.Unconscious) {
			return
		}
		reset()
	})

	fx.On(func(s *core.AfterTakeDamageState) {
		if s.Source.HitPoints > 0 {
			return
		}
		if !s.Source.Conditions.Match(tags.Unconscious) {
			s.Source.RemoveCondition(tags.Stable, nil)
			s.Source.AddCondition(tags.Unconscious, fx)
			reset()
			return
		}
		amount := 1
		if s.Result.IsCriticalSuccess() {
			amount = 2
		}
		failures += amount
		s.Source.Dispatcher.Begin(core.DeathSavingThrowAutomaticEvent{Source: s.Source, Failure: true})
		defer s.Source.Dispatcher.End()
		s.Source.Dispatcher.Emit(core.DeathSavingThrowResultEvent{Source: s.Source, Success: success, Failure: failures})
		if checkStatus(s.Source) && s.Result.Value > s.Source.MaxHitPoints {
			s.Source.Die()
		}
	})

	fx.On(func(s *core.TurnState) {
		if !s.Source.MatchCondition(tags.Unconscious) {
			return
		}
		s.Source.Dispatcher.Begin(core.DeathSavingThrowEvent{Source: s.Source})
		defer s.Source.Dispatcher.End()
		result := s.Source.SaveThrow(tags.HitPoints, 10)
		if result.Success {
			success++
			if result.Critical {
				reset()
				s.Source.Dispatcher.Begin(core.DeathSavingThrowAutomaticEvent{Source: s.Source, Failure: false})
				defer s.Source.Dispatcher.End()
				s.Source.RemoveCondition(tags.Unconscious, nil)
				s.Source.ModifyAttribute(tags.HitPoints, 1, "Death Saving Throw critical success")
				return
			}
		} else {
			failures++
			if result.Critical {
				failures++
			}
		}
		s.Source.Dispatcher.Emit(core.DeathSavingThrowResultEvent{Source: s.Source, Success: success, Failure: failures})
		checkStatus(s.Source)
	})

	return fx
}
