package base

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
	"anvil/internal/tag"
)

//nolint:gocognit
func NewAttackOfOpportunityEffect() *core.Effect {
	fx := &core.Effect{Name: "Attack Of Opportunity"}

	fx.WithBeforeMoveStep(func(_ *core.Effect, s *core.MoveState) {
		if s.Action != nil && s.Action.Tags().MatchTag(tags.Teleport) {
			return
		}
		enemies := s.Source.World.ActorsInRange(
			s.From,
			1,
			func(other *core.Actor) bool { return other.IsHostileTo(s.Source) },
		)
		options := []core.RequestOption{
			{Value: true, Label: "Yes", Default: true},
			{Value: false, Label: "No"},
		}
		for _, other := range enemies {
			if !other.CanAct() {
				continue
			}
			if s.Source.Encounter.IsOver() {
				return
			}
			if !other.Resources.CanAfford(map[tag.Tag]int{tags.Reaction: 1}) {
				continue
			}
			response := s.Source.World.Ask(other, "Take attack of opportunity?", options)
			b, ok := response.Value.(bool)
			if !ok || !b {
				continue
			}
			baseAttack := other.BestWeaponAttack()
			if baseAttack == nil {
				continue
			}
			s.Source.Dispatcher.Start(core.EffectType, core.EffectEvent{Source: s.Source, Effect: fx})
			other.ConsumeResource(tags.Reaction, 1)
			baseAttack.Perform([]grid.Position{s.Source.Position}, false)
			s.Source.Dispatcher.End()
		}
	})
	return fx
}
