package fighter

import (
	"anvil/internal/core"
	"anvil/internal/core/tags"
	"anvil/internal/tag"
)

func NewFightingStyleDefense() *core.Effect {
	fx := &core.Effect{Name: "Fighting Style: Defense"}
	fx.WithAttributeCalculation(func(_ *core.Effect, s *core.AttributeCalculationState) {
		if !s.Attribute.MatchExact(tags.ArmorClass) {
			return
		}
		valid := tag.ContainerFromTag(tags.LightArmor, tags.MediumArmor, tags.MediumArmor, tags.Shield)
		trigger := false
		for _, e := range s.Source.Equipped {
			if e.Tags().HasAnyTag(valid) {
				trigger = true
			}
		}
		if !trigger {
			return
		}
		s.Expression.AddScalar(1, fx.Name)
	})
	return fx
}
