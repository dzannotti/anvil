package basic

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/grid"
	"anvil/internal/tag"
)

//nolint:gocognit
func NewAttackOfOpportunityEffect() *core.Effect {
	fx := &core.Effect{Name: "Attack Of Opportunity"}

	fx.On(func(s *core.PreMoveStep) {
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

			if !other.Resources.CanAfford(map[tag.Tag]int{tags.ResourceReaction: 1}) {
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

			s.Source.Dispatcher.Begin(core.EffectEvent{Source: s.Source, Effect: fx})
			other.ConsumeResource(tags.ResourceReaction, 1)
			// TODO: Create proper AOO action with Reaction cost instead of Action cost
			baseAttack.Perform([]grid.Position{s.Source.Position})
			s.Source.Dispatcher.End()
		}
	})
	return fx
}
