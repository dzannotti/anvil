package base

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
)

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

	fx.WithSerialize(func(_ *core.Effect, s *core.SerializeState) {
		s.State.Data = map[string]int{
			"Success":  success,
			"Failures": failures,
		}
	})

	fx.WithDeserialize(func(_ *core.Effect, s *core.SerializeState) {
		data, ok := s.State.Data.(map[string]interface{})
		if !ok {
			panic("could not deserialize dst")
		}
		success = int(data["Success"].(float64))
		failures = int(data["Failures"].(float64))
	})

	fx.WithAttributeChanged(func(_ *core.Effect, s *core.AttributeChangedState) {
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

	fx.WithConditionRemoved(func(_ *core.Effect, s *core.ConditionChangedState) {
		if !s.Condition.Match(tags.Unconscious) {
			return
		}
		reset()
	})

	fx.WithAfterTakeDamage(func(_ *core.Effect, s *core.AfterTakeDamageState) {
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
		failures = failures + amount
		s.Source.Log.Start(
			core.DeathSavingThrowAutomaticType,
			core.DeathSavingThrowAutomaticEvent{Source: s.Source, Failure: true},
		)
		defer s.Source.Log.End()
		s.Source.Log.Add(
			core.DeathSavingThrowResultType,
			core.DeathSavingThrowResultEvent{Source: s.Source, Success: success, Failure: failures},
		)
		if checkStatus(s.Source) && s.Result.Value > s.Source.MaxHitPoints {
			s.Source.Die()
		}
	})

	fx.WithTurnStarted(func(_ *core.Effect, s *core.TurnState) {
		if !s.Source.MatchCondition(tags.Unconscious) {
			return
		}
		s.Source.Log.Start(core.DeathSavingThrowType, core.DeathSavingThrowEvent{Source: s.Source})
		defer s.Source.Log.End()
		result := s.Source.SaveThrow(tags.HitPoints, 10)
		if result.Success {
			success = success + 1
			if result.Critical {
				reset()
				s.Source.Log.Start(
					core.DeathSavingThrowAutomaticType,
					core.DeathSavingThrowAutomaticEvent{Source: s.Source, Failure: false},
				)
				defer s.Source.Log.End()
				s.Source.RemoveCondition(tags.Unconscious, nil)
				s.Source.ModifyAttribute(tags.HitPoints, 1, "Death Saving Throw critical success")
				return
			}
		} else {
			failures = failures + 1
			if result.Critical {
				failures = failures + 1
			}
		}
		s.Source.Log.Add(
			core.DeathSavingThrowResultType,
			core.DeathSavingThrowResultEvent{Source: s.Source, Success: success, Failure: failures},
		)
		checkStatus(s.Source)
	})

	return fx
}
