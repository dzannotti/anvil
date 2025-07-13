package fighter

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/tag"

	"github.com/google/uuid"
)

func NewFightingStyleDefense() *core.Effect {
	fx := &core.Effect{
		Archetype: "fighting-style-defense",
		ID:        uuid.New().String(),
		Name:      "Fighting Style: Defense",
	}
	fx.On(func(s *core.AttributeCalculation) {
		if !s.Attribute.MatchExact(tags.ActorArmorClass) {
			return
		}

		valid := tag.NewContainer(tags.LightArmor, tags.MediumArmor, tags.HeavyArmor, tags.Shield)
		trigger := false
		for _, e := range s.Source.Equipped {
			if e.Tags().HasAny(valid) {
				trigger = true
			}
		}

		if !trigger {
			return
		}

		s.Expression.AddConstant(1, fx.Name)
	})
	return fx
}
